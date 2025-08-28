package service

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/internal/repository/db"
	"Art-Design-Backend/pkg/ai"
	"Art-Design-Backend/pkg/aliyun"
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/slicer_client"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/pgvector/pgvector-go"
	"go.uber.org/zap"
)

type KnowledgeBaseService struct {
	OssClient         *aliyun.OssClient             // 阿里云OSS
	Slicer            *slicer_client.Slicer         // 文档切片
	GormTX            *db.GormTransactionManager    // 事务
	AIModelClient     *ai.AIModelClient             // AI 模型
	KnowledgeBaseRepo *repository.KnowledgeBaseRepo // 知识库
	AIProviderRepo    *repository.AIProviderRepo    // AI 供应商
	UserRepo          *repository.UserRepo          //  用户
}

func (k *KnowledgeBaseService) UploadAndVectorizeDocument(
	c context.Context,
	file multipart.File,
	filename string,
	fileSize int64,
) error {
	// Step 1: 上传文档到 OSS
	documentURL, err := k.OssClient.UploadKnowledgeBaseFile(c, filename, file)
	if err != nil {
		zap.L().Error("上传文档失败", zap.Error(err))
		return fmt.Errorf("上传文档失败: %w", err)
	}

	// Step 2: 创建 KnowledgeBaseFile 记录
	knowledgeBaseFile := &entity.KnowledgeBaseFile{
		OriginalFileName: strings.TrimSuffix(filename, filepath.Ext(filename)), // 去掉文件后缀（如将 "example.pdf" 变为 "example"）
		FileType:         strings.ToLower(filepath.Ext(filename)[1:]),          // 去掉 "." => pdf/docx/txt
		FileSize:         fileSize,                                             // 单位：字节
		FilePath:         documentURL,                                          // OSS 返回的文件存储路径或 URL
	}
	if err = k.KnowledgeBaseRepo.CreateKnowledgeFile(c, knowledgeBaseFile); err != nil {
		zap.L().Error("保存 知识库文件 失败", zap.Error(err))
		return fmt.Errorf("保存 知识库文件 失败: %w", err)
	}

	// Step 3: 文档分块
	chunks, err := k.Slicer.GetChunksFromSlicer(documentURL)
	if err != nil {
		zap.L().Error("文档分块失败", zap.Error(err))
		return fmt.Errorf("文档分块失败: %w", err)
	}
	if len(chunks) == 0 {
		return fmt.Errorf("文档内容为空，无法切分")
	}

	// 使用事务包裹整个处理流程
	err = k.GormTX.Transaction(c, func(ctx context.Context) error {
		// Step 4: 保存分块内容
		chunkList := make([]*entity.FileChunk, 0, len(chunks))
		for i, chunk := range chunks {
			chunkEntity := &entity.FileChunk{
				FileID:     knowledgeBaseFile.ID,
				ChunkIndex: i,
				Content:    chunk,
			}
			if err = k.KnowledgeBaseRepo.CreateFileChunk(ctx, chunkEntity); err != nil {
				zap.L().Error("创建文件块失败", zap.Int("index", i), zap.Error(err))
				return fmt.Errorf("创建文件块失败(index %d): %w", i, err)
			}
			chunkList = append(chunkList, chunkEntity)
		}

		// Step 5: 分批获取 Embedding（每次最多 10 个 chunk）
		batchSize := 10
		allEmbeddings := make([][]float32, 0, len(chunkList))
		for i := 0; i < len(chunks); i += batchSize {
			end := i + batchSize
			if end > len(chunks) {
				end = len(chunks)
			}
			batchChunks := chunks[i:end]

			// 调用千问 Embedding API（每次最多 10 个）
			batchEmbeddings, err := k.getQianwenEmbeddings(ctx, batchChunks)
			if err != nil {
				return fmt.Errorf("获取 Embedding 失败(batch %d-%d): %w", i, end-1, err)
			}
			allEmbeddings = append(allEmbeddings, batchEmbeddings...)
		}

		// Step 6: 检查向量数量是否匹配
		if len(allEmbeddings) != len(chunkList) {
			zap.L().Error("向量数量与分块数量不一致",
				zap.Int("chunks", len(chunkList)),
				zap.Int("vectors", len(allEmbeddings)))
			return fmt.Errorf("向量数量与分块数量不一致")
		}

		// Step 7: 保存向量
		for i, chunk := range chunkList {
			chunkVector := &entity.ChunkVector{
				ChunkID:   chunk.ID,
				Embedding: pgvector.NewVector(allEmbeddings[i]),
			}
			if err = k.KnowledgeBaseRepo.CreateChunkVector(ctx, chunkVector); err != nil {
				zap.L().Error("保存向量失败", zap.Int64("chunkID", chunk.ID), zap.Error(err))
				return fmt.Errorf("保存向量失败(chunkID %d): %w", chunk.ID, err)
			}
		}

		// ✅ 输出日志（事务内）
		zap.L().Info("文档上传与向量化完成",
			zap.Int64("fileID", knowledgeBaseFile.ID),
			zap.Int("chunkCount", len(chunks)),
			zap.String("file", filename),
		)
		return nil
	})

	if err != nil {
		// 事务已自动回滚，此处可补充额外日志
		zap.L().Error("文档处理事务失败", zap.Error(err))
		return err
	}

	return nil
}

// 获取嵌入向量
func (k *KnowledgeBaseService) getQianwenEmbeddings(c context.Context, chunks []string) ([][]float32, error) {
	const providerID int64 = 51088793876300041

	provider, err := k.AIProviderRepo.GetAIProviderByIDWithCache(c, providerID)
	if err != nil {
		zap.L().Error("获取嵌入模型供应商失败", zap.Error(err))
		return nil, fmt.Errorf("获取嵌入模型供应商失败: %w", err)
	}

	embeddings, err := k.AIModelClient.Embed(c, provider.APIKey, chunks)
	if err != nil {
		zap.L().Error("获取嵌入向量失败", zap.Error(err))
		return nil, fmt.Errorf("获取嵌入向量失败: %w", err)
	}

	return embeddings, nil
}

func (k *KnowledgeBaseService) GetKnowledgeBaseFileList(
	ctx context.Context,
	req *query.KnowledgeBaseFile,
) (*common.PaginationResp[response.KnowledgeBaseFile], error) {

	// 1. 查用户 ID（根据 CreateUser 模糊匹配用户名）
	var filterUserIDs []int64
	var err error
	if req.CreateUser != nil {
		filterUserIDs, err = k.UserRepo.GetUserIDsByName(ctx, *req.CreateUser)
		if err != nil {
			return nil, err
		}
		if len(filterUserIDs) == 0 {
			// zap.L().Info("未找到用户") // 可以记录日志，但不是错误
			return common.BuildPageResp[response.KnowledgeBaseFile]([]response.KnowledgeBaseFile{}, 0, req.PaginationReq), nil
		}
	}

	// 2. 分页查询文件实体
	fileEntities, total, err := k.KnowledgeBaseRepo.GetKnowledgeFilePage(ctx, req, filterUserIDs)
	if err != nil {
		return nil, err
	}

	// 3. 收集所有上传者 ID
	uploaderIDs := make([]int64, 0, len(fileEntities))
	for _, f := range fileEntities {
		uploaderIDs = append(uploaderIDs, f.CreateBy)
	}

	// 4. 批量查询用户信息（ID -> 用户名）
	userMap, err := k.UserRepo.GetUserMapByIDs(ctx, uploaderIDs)
	if err != nil {
		return nil, err
	}

	// 5. 转换为响应对象
	fileResponses := make([]response.KnowledgeBaseFile, 0, len(fileEntities))
	for _, fileEntity := range fileEntities {
		var fileResp response.KnowledgeBaseFile
		_ = copier.Copy(&fileResp, fileEntity)

		// 替换成用户名
		if username, ok := userMap[fileEntity.CreateBy]; ok {
			fileResp.CreateUser = username
		}

		fileResponses = append(fileResponses, fileResp)
	}

	// 6. 构建分页响应
	resp := common.BuildPageResp[response.KnowledgeBaseFile](fileResponses, total, req.PaginationReq)
	return resp, nil
}

func (k *KnowledgeBaseService) GetKnowledgeBasePage(c context.Context, q *query.KnowledgeBase) (res []*response.KnowledgeBase, err error) {
	knowledgeBases, err := k.KnowledgeBaseRepo.GetKnowledgeBasePage(c, q, authutils.GetUserID(c))
	res = make([]*response.KnowledgeBase, 0, len(knowledgeBases))
	for _, knowledgeBase := range knowledgeBases {
		var knowledgeBaseResp response.KnowledgeBase
		_ = copier.Copy(&knowledgeBaseResp, knowledgeBase)
		res = append(res, &knowledgeBaseResp)
	}
	return
}

func (k *KnowledgeBaseService) CreateKnowledgeBase(c context.Context, req *request.KnowledgeBase) (err error) {
	var knowledgeBase entity.KnowledgeBase
	_ = copier.Copy(&knowledgeBase, req)
	err = k.KnowledgeBaseRepo.CreateKnowledgeBase(c, &knowledgeBase)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return
	}
	var knowledgeBaseFileRelList []*entity.KnowledgeBaseFileRel
	knowledgeBaseFileRelList = make([]*entity.KnowledgeBaseFileRel, 0, len(req.Files))
	for _, file := range req.Files {
		knowledgeBaseFileRelList = append(knowledgeBaseFileRelList, &entity.KnowledgeBaseFileRel{
			KnowledgeBaseID:     knowledgeBase.ID,
			KnowledgeBaseFileID: int64(file.ID),
		})
	}
	if err = k.KnowledgeBaseRepo.CreateKnowledgeBaseFileRel(c, knowledgeBaseFileRelList); err != nil {
		zap.L().Error("创建知识库文件关系失败", zap.Error(err))
		return
	}
	return
}

func (k *KnowledgeBaseService) DeleteKnowledgeBase(c *gin.Context, id int64) (err error) {
	err = k.GormTX.Transaction(c, func(ctx context.Context) (err error) {
		if err = k.KnowledgeBaseRepo.DeleteKnowledgeBase(c, id); err != nil {
			zap.L().Error("删除知识库失败", zap.Error(err))
			return
		}
		if err = k.KnowledgeBaseRepo.DeleteKnowledgeBaseFileRel(c, id); err != nil {
			zap.L().Error("删除知识库文件关系失败", zap.Error(err))
			return
		}
		return
	})
	if err != nil {
		zap.L().Error("删除知识库事务失败", zap.Error(err))
		return
	}
	return
}

func (k *KnowledgeBaseService) UpdateKnowledgeBase(c context.Context, r *request.KnowledgeBase) (err error) {
	var knowledgeBase entity.KnowledgeBase
	_ = copier.Copy(&knowledgeBase, r)
	err = k.GormTX.Transaction(c, func(ctx context.Context) (err error) {
		if err = k.KnowledgeBaseRepo.DeleteKnowledgeBaseFileRel(ctx, knowledgeBase.ID); err != nil {
			zap.L().Error("删除知识库文件关系失败", zap.Error(err))
			return
		}
		// 创建知识库文件关系,只有文件不为空时才创建
		if len(r.Files) != 0 {
			// 预分配足够容量的切片，避免多次扩容
			knowledgeBaseFileRelList := make([]*entity.KnowledgeBaseFileRel, 0, len(r.Files))

			for _, file := range r.Files {
				knowledgeBaseFileRelList = append(knowledgeBaseFileRelList, &entity.KnowledgeBaseFileRel{
					KnowledgeBaseID:     knowledgeBase.ID,
					KnowledgeBaseFileID: int64(file.ID),
				})
			}
			if err = k.KnowledgeBaseRepo.CreateKnowledgeBaseFileRel(ctx, knowledgeBaseFileRelList); err != nil {
				zap.L().Error("创建知识库文件关系失败", zap.Error(err))
				return
			}
		}
		if err = k.KnowledgeBaseRepo.UpdateKnowledgeBase(ctx, &knowledgeBase); err != nil {
			zap.L().Error("更新知识库失败", zap.Error(err))
			return
		}
		return
	})
	if err != nil {
		zap.L().Error("更新知识库事务失败", zap.Error(err))
		return
	}
	return
}

func (k *KnowledgeBaseService) GetKnowledgeBaseFilesByID(
	c context.Context, id int64) (res []*response.KnowledgeBaseFile, err error) {
	knowledgeBaseFiles, err := k.KnowledgeBaseRepo.GetKnowledgeBaseFilesByID(c, id)
	res = make([]*response.KnowledgeBaseFile, 0, len(knowledgeBaseFiles))
	for _, knowledgeBaseFile := range knowledgeBaseFiles {
		var knowledgeBaseFileResp response.KnowledgeBaseFile
		_ = copier.Copy(&knowledgeBaseFileResp, knowledgeBaseFile)
		res = append(res, &knowledgeBaseFileResp)
	}
	return
}

func (k *KnowledgeBaseService) GetSimpleKnowledgeBaseList(c context.Context) (res []*response.SimpleKnowledgeBase, err error) {
	var userID int64
	userID = authutils.GetUserID(c)
	knowledgeBases, err := k.KnowledgeBaseRepo.GetSimpleKnowledgeBaseList(c, userID)
	res = make([]*response.SimpleKnowledgeBase, 0, len(knowledgeBases))
	for _, knowledgeBase := range knowledgeBases {
		var knowledgeBaseResp response.SimpleKnowledgeBase
		_ = copier.Copy(&knowledgeBaseResp, knowledgeBase)
		res = append(res, &knowledgeBaseResp)
	}
	return
}

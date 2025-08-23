package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/internal/repository/db"
	"Art-Design-Backend/pkg/ai"
	"Art-Design-Backend/pkg/aliyun"
	"Art-Design-Backend/pkg/slicer_client"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

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
		OriginalFileName: filename,                                    // 用户上传的原始文件名
		FileType:         strings.ToLower(filepath.Ext(filename)[1:]), // 去掉 "." => pdf/docx/txt
		FileSize:         fileSize,                                    // 单位：字节
		FilePath:         documentURL,                                 // OSS 返回的文件存储路径或 URL
	}
	if err = k.KnowledgeBaseRepo.CreateAgentFile(c, knowledgeBaseFile); err != nil {
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

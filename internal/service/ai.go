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
	"Art-Design-Backend/pkg/slicer_client"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/pgvector/pgvector-go"
	"go.uber.org/zap"
	"mime/multipart"
	"strings"
)

type AIService struct {
	AIModelRepo    *repository.AIModelRepo    // 模型Repo
	AIModelClient  *ai.AIModelClient          // 聊天
	AIProviderRepo *repository.AIProviderRepo // 模型供应商Repo
	AIAgentRepo    *repository.AIAgentRepo    // 智能助手Repo
	OssClient      *aliyun.OssClient          // 阿里云OSS
	Slicer         *slicer_client.Slicer      // 文档切片
	GormTX         *db.GormTransactionManager // 事务
}

func (a *AIService) CreateAIProvider(c context.Context, r *request.AIProvider) (err error) {
	var aiProvider entity.AIProvider
	_ = copier.Copy(&aiProvider, &r)
	if err = a.AIProviderRepo.CheckAIDuplicate(c, &aiProvider); err != nil {
		zap.L().Error(err.Error())
		return
	}
	if err = a.AIProviderRepo.Create(c, &aiProvider); err != nil {
		zap.L().Error(err.Error())
	}
	return
}

func (a *AIService) GetAIProviderPage(c context.Context, q *query.AIProvider) (res *common.PaginationResp[*response.AIProvider], err error) {
	var aiProviders []*entity.AIProvider
	var total int64
	if aiProviders, total, err = a.AIProviderRepo.GetAIProviderPage(c, q); err != nil {
		zap.L().Error(err.Error())
		return
	}
	aiProviderRes := make([]*response.AIProvider, 0, len(aiProviders))
	for _, provider := range aiProviders {
		var aiProviderResp response.AIProvider
		_ = copier.Copy(&aiProviderResp, provider)
		aiProviderRes = append(aiProviderRes, &aiProviderResp)
	}
	res = common.BuildPageResp[*response.AIProvider](aiProviderRes, total, q.PaginationReq)
	return
}

func (a *AIService) CreateAIModel(c context.Context, r *request.AIModel) (err error) {
	var aiModel entity.AIModel
	_ = copier.Copy(&aiModel, &r)
	if err = a.AIModelRepo.CheckAIDuplicate(c, &aiModel); err != nil {
		zap.L().Error("AI模型已存在", zap.Error(err))
		return
	}
	if err = a.AIModelRepo.Create(c, &aiModel); err != nil {
		zap.L().Error("aiModelCreate失败", zap.Error(err))
		return
	}
	return
}
func (a *AIService) GetSimpleChatModelList(c context.Context) (res []*response.SimpleAIModel, err error) {
	res, err = a.AIModelRepo.GetSimpleChatModelList(c)
	if err != nil {
		zap.L().Error("GetSimpleChatModelList 失败", zap.Error(err))
		return
	}
	return
}

func (a *AIService) GetAIModelPage(c context.Context, q *query.AIModel) (res *common.PaginationResp[*response.AIModel], err error) {
	// 查询模型分页数据
	modelList, total, err := a.AIModelRepo.GetAIModelPage(c, q)
	if err != nil {
		zap.L().Error("AI模型分页查询失败", zap.Error(err))
		return
	}

	// 去重 ProviderID
	providerIDSet := make(map[int64]struct{})
	for _, model := range modelList {
		providerIDSet[model.ProviderID] = struct{}{}
	}

	// 转换为 ID 列表
	providerIDs := make([]int64, 0, len(providerIDSet))
	for id := range providerIDSet {
		providerIDs = append(providerIDs, id)
	}

	// 查询 Provider 列表
	providerList, err := a.AIProviderRepo.GetProviderNameByIDList(c, providerIDs)
	if err != nil {
		zap.L().Error("GetProviderNameByIDList 失败", zap.Error(err))
		return
	}

	// 构建 ID → Name 映射
	providerNameMap := make(map[int64]string, len(providerList))
	for _, p := range providerList {
		providerNameMap[p.ID] = p.Name
	}

	// 构造响应列表
	resList := make([]*response.AIModel, 0, len(modelList))
	for _, model := range modelList {
		resp := &response.AIModel{}
		_ = copier.Copy(resp, model)
		resp.Provider = providerNameMap[model.ProviderID]
		resList = append(resList, resp)
	}

	// 返回分页响应
	res = common.BuildPageResp[*response.AIModel](resList, total, q.PaginationReq)
	return
}

func (a *AIService) GetSimpleProviderList(c context.Context) (res []*response.SimpleAIProvider, err error) {
	var aiProviders []*entity.AIProvider
	aiProviders, err = a.AIProviderRepo.GetSimpleProviderList(c)
	if err != nil {
		zap.L().Error("db GetSimpleProviderList", zap.Error(err))
	}
	res = make([]*response.SimpleAIProvider, 0, len(aiProviders))
	for _, provider := range aiProviders {
		var aiProviderResp response.SimpleAIProvider
		_ = copier.Copy(&aiProviderResp, provider)
		res = append(res, &aiProviderResp)
	}
	return
}

func (a *AIService) UploadModelIcon(c *gin.Context, filename string, src multipart.File) (url string, err error) {
	url, err = a.OssClient.UploadModelIcon(c, filename, src)
	if err != nil {
		zap.L().Error("上传模型图标失败", zap.Error(err))
		return
	}
	return
}

// ChatCompletion
// todo
// 1. 对话管理
// 2. 模型信息缓存 ✅
// 3. 标签抽取
// 4. 价格计算
func (a *AIService) ChatCompletion(c *gin.Context, r *request.ChatCompletion) (err error) {
	modelID := int64(r.ID)
	modelInfo, err := a.AIModelRepo.GetAIModelByID(c, modelID)
	if err != nil {
		zap.L().Error("获取AI模型失败", zap.Int64("model_id", modelID), zap.Error(err))
		return
	}
	provider, err := a.AIProviderRepo.GetAIProviderByIDWithCache(c, modelInfo.ProviderID)
	if err != nil {
		zap.L().Error("获取AI模型供应商失败", zap.Int64("provider_id", modelInfo.ProviderID), zap.Error(err))
		return
	}
	// 直接调用新版 ChatStreamWithWriter
	err = a.AIModelClient.ChatStreamWithWriter(
		c.Request.Context(), c.Writer,
		provider.BaseURL+modelInfo.ApiPath,
		provider.APIKey,
		ai.DefaultStreamChatRequest(modelInfo.Model, r.Messages),
	)
	if err != nil {
		zap.L().Error("AI模型聊天失败", zap.Error(err))
		return
	}
	return
}

func (a *AIService) getQianwenEmbeddings(c context.Context, chunks []string) ([][]float32, error) {
	const providerID int64 = 51088793876300041

	provider, err := a.AIProviderRepo.GetAIProviderByIDWithCache(c, providerID)
	if err != nil {
		zap.L().Error("获取嵌入模型供应商失败", zap.Error(err))
		return nil, fmt.Errorf("获取嵌入模型供应商失败: %w", err)
	}

	embeddings, err := a.AIModelClient.Embed(c, provider.APIKey, chunks)
	if err != nil {
		zap.L().Error("获取嵌入向量失败", zap.Error(err))
		return nil, fmt.Errorf("获取嵌入向量失败: %w", err)
	}

	return embeddings, nil
}

func (a *AIService) UploadAndVectorizeDocument(
	c context.Context,
	file multipart.File,
	filename string,
	agentID int64,
) error {
	// Step 1: 上传文档到 OSS
	documentURL, err := a.OssClient.UploadAgentDocument(c, filename, file)
	if err != nil {
		zap.L().Error("上传文档失败", zap.Error(err))
		return fmt.Errorf("上传文档失败: %w", err)
	}

	// Step 2: 创建 AgentFile 记录
	agentFile := &entity.AgentFile{
		AgentID: agentID,
		FileURL: documentURL,
	}
	if err = a.AIAgentRepo.CreateAgentFile(c, agentFile); err != nil {
		zap.L().Error("保存 AgentFile 失败", zap.Error(err))
		return fmt.Errorf("保存 AgentFile 失败: %w", err)
	}

	// Step 3: 文档分块
	chunks, err := a.Slicer.GetChunksFromSlicer(documentURL)
	if err != nil {
		zap.L().Error("文档分块失败", zap.Error(err))
		return fmt.Errorf("文档分块失败: %w", err)
	}
	if len(chunks) == 0 {
		return fmt.Errorf("文档内容为空，无法切分")
	}

	// Step 4: 保存分块内容
	chunkList := make([]*entity.FileChunk, 0, len(chunks))
	for i, chunk := range chunks {
		chunkEntity := &entity.FileChunk{
			FileID:     agentFile.ID,
			ChunkIndex: i,
			Content:    chunk,
		}
		if err = a.AIAgentRepo.CreateFileChunk(c, chunkEntity); err != nil {
			zap.L().Error("创建文件块失败", zap.Int("index", i), zap.Error(err))
			return fmt.Errorf("创建文件块失败(index %d): %w", i, err)
		}
		chunkList = append(chunkList, chunkEntity)
	}

	// Step 5: 获取 Embedding（使用指定 provider,千问）
	embeddings, err := a.getQianwenEmbeddings(c, chunks)
	if err != nil {
		return err
	}
	if len(embeddings) != len(chunkList) {
		zap.L().Error("向量数量与分块数量不一致",
			zap.Int("chunks", len(chunkList)),
			zap.Int("vectors", len(embeddings)))
		return fmt.Errorf("向量数量与分块数量不一致")
	}

	// Step 6: 保存向量
	for i, chunk := range chunkList {
		chunkVector := &entity.ChunkVector{
			ChunkID:   chunk.ID,
			Embedding: pgvector.NewVector(embeddings[i]),
		}
		if err = a.AIAgentRepo.CreateChunkVector(c, chunkVector); err != nil {
			zap.L().Error("保存向量失败", zap.Int64("chunkID", chunk.ID), zap.Error(err))
			return fmt.Errorf("保存向量失败(chunkID %d): %w", chunk.ID, err)
		}
	}

	// ✅ 输出日志
	zap.L().Info("文档上传与向量化完成",
		zap.Int64("agentID", agentID),
		zap.Int("chunkCount", len(chunks)),
		zap.String("file", filename),
	)
	return nil
}

func (a *AIService) CreateAgent(c *gin.Context, r *request.AIAgent) (err error) {
	agent := &entity.AIAgent{}
	_ = copier.Copy(agent, r)
	if err = a.AIAgentRepo.Create(c, agent); err != nil {
		zap.L().Error("创建AI模型失败", zap.Error(err))
		return
	}
	return
}

func (a *AIService) GetAIAgentPage(c *gin.Context, q *query.AIAgent) (res *common.PaginationResp[*response.AIAgent], err error) {
	agents, total, err := a.AIAgentRepo.GetAIAgentPage(c, q)
	if err != nil {
		zap.L().Error("获取AI模型分页失败", zap.Error(err))
		return
	}
	agentRes := make([]*response.AIAgent, 0, len(agents))
	for _, agent := range agents {
		agentResp := &response.AIAgent{}
		_ = copier.Copy(agentResp, agent)
		agentRes = append(agentRes, agentResp)
	}
	res = common.BuildPageResp[*response.AIAgent](agentRes, total, q.PaginationReq)
	return
}

func (a *AIService) GetSimpleAgentList(c context.Context) (res []*response.SimpleAIAgent, err error) {
	var agentList []*entity.AIAgent
	agentList, err = a.AIAgentRepo.GetSimpleAgentList(c)
	if err != nil {
		zap.L().Error("获取智能体列表失败", zap.Error(err))
		return
	}
	res = make([]*response.SimpleAIAgent, 0, len(agentList))
	for _, agent := range agentList {
		agentResp := &response.SimpleAIAgent{}
		_ = copier.Copy(agentResp, agent)
		res = append(res, agentResp)
	}
	return
}

func (a *AIService) AgentChatCompletion(c *gin.Context, r *request.ChatCompletion) (err error) {
	// Step 1: 获取 Agent 信息
	agentInfo, err := a.AIAgentRepo.GetAIAgentByIDWithCache(c, int64(r.ID))
	if err != nil {
		zap.L().Error("获取智能体失败", zap.Error(err))
		return fmt.Errorf("获取智能体失败: %w", err)
	}

	// Step 2: 获取用户最新提问（最后一条 user 消息）
	var latestQuestion string
	for i := len(r.Messages) - 1; i >= 0; i-- {
		if r.Messages[i].Role == "user" {
			latestQuestion = r.Messages[i].Content
			break
		}
	}
	if latestQuestion == "" {
		return fmt.Errorf("未找到用户提问")
	}

	// Step 3: 向量化提问
	embedding, err := a.getQianwenEmbeddings(c, []string{latestQuestion})
	if err != nil {
		return err
	}

	// Step 4: 基于 embedding 搜索相关内容（假设你有一个向量搜索接口）
	retrievedTexts, err := a.AIAgentRepo.SearchAgentRelatedChunks(c, agentInfo.ID, embedding[0])
	if err != nil {
		zap.L().Error("向量检索失败", zap.Error(err))
		return fmt.Errorf("向量检索失败: %w", err)
	}

	// Step 5: 构造新的对话上下文（system + retrieved + user）
	var fullMessages []ai.ChatMessage

	// 添加 agent 的默认 system prompt（如有）
	if agentInfo.SystemPrompt != "" {
		fullMessages = append(fullMessages, ai.ChatMessage{
			Role:    "system",
			Content: agentInfo.SystemPrompt,
		})
	}

	// 将检索到的文本合并为一段 context，并构造辅助 system 提示
	if len(retrievedTexts) > 0 {
		contextText := strings.Join(retrievedTexts, "\n\n")
		fullMessages = append(fullMessages, ai.ChatMessage{
			Role: "system",
			Content: fmt.Sprintf(`以下是与用户问题相关的背景资料。请仅根据这些资料回答问题，如果资料中没有提及，请直接回复：“很抱歉，我无法在现有知识中找到相关答案。”：

%s`, contextText),
		})
	} else {
		fullMessages = append(fullMessages, ai.ChatMessage{
			Role:    "system",
			Content: `注意：当前没有任何与用户问题相关的背景资料。如果资料中没有提及，请直接回复：“很抱歉，我无法在现有知识中找到相关答案。”`,
		})
	}

	// 拼接用户原始对话
	fullMessages = append(fullMessages, r.Messages...)

	// Step 6: 调用模型回答（根据 agentInfo.ModelID 决定使用哪个模型）
	return a.ChatCompletion(c, &request.ChatCompletion{
		ID:       common.LongStringID(agentInfo.ModelID),
		Messages: fullMessages,
	})
}

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
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type AIService struct {
	AIModelRepo       *repository.AIModelRepo       // 模型Repo
	AIModelClient     *ai.AIModelClient             // 聊天
	AIProviderRepo    *repository.AIProviderRepo    // 模型供应商Repo
	KnowledgeBaseRepo *repository.KnowledgeBaseRepo // 知识库Repo
	ConversationRepo  *repository.ConversationRepo  // 会话Repo
	OssClient         *aliyun.OssClient             // 阿里云OSS
	Slicer            *slicer_client.Slicer         // 文档切片
	GormTX            *db.GormTransactionManager    // 事务
}

// 获取嵌入向量
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
// 1. 对话管理 ✅
// 2. 模型信息缓存 ✅
// 3. 标签抽取
// 4. 价格计算
func (a *AIService) ChatCompletion(c *gin.Context, r *request.ChatCompletion) (err error) {
	modelID := int64(r.ID)

	var conversation entity.Conversation

	// Step 1: 获取模型信息
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

	// Step 2: 提取用户最新问题
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

	// Step 3: 会话初始化
	if r.ConversationID == 0 {
		// Step 3.1: 创建会话
		conversation = entity.Conversation{}
		// 确保字符少于等于10时才进行截取，否则由大模型总结为十个字
		if utf8.RuneCountInString(latestQuestion) <= 10 {
			conversation.Title = latestQuestion
		} else {
			// Step 3.2: 如果新对话提问信息过长，使用当前选择的对话大模型总结十个字作为会话标题
			titleSummaryJson, err := a.AIModelClient.ChatRequest(
				c.Request.Context(),
				provider.BaseURL+modelInfo.ApiPath,
				provider.APIKey,
				ai.DefaultChatRequest(modelInfo.Model,
					[]ai.ChatMessage{
						{
							Role:    "system",
							Content: ai.TitleSummaryPrompt,
						},
						{
							Role:    "user",
							Content: latestQuestion,
						},
					},
				),
			)
			if err != nil {
				zap.L().Error("总结标题失败", zap.Error(err))
				conversation.Title = latestQuestion[:10]
			} else {
				var titleSummary ai.ChatCompletionResponse
				_ = sonic.Unmarshal(titleSummaryJson, &titleSummary)
				conversation.Title = titleSummary.FirstText()
				zap.L().Info("总结标题成功", zap.Int64("conversation_id", conversation.ID), zap.String("title", conversation.Title))
			}
		}

		if err = a.ConversationRepo.CreateConversation(c, &conversation); err != nil {
			zap.L().Error("创建会话失败", zap.Error(err))
			return
		}
	}

	// Step 4: 构建知识库回答
	var fullMessages []ai.ChatMessage
	var documents []string
	var retrievedTextChunks []*entity.FileChunk
	if r.KnowledgeBaseID != 0 {
		// Step 3: 计算 embedding
		embedding, err := a.getQianwenEmbeddings(c, []string{latestQuestion})
		if err != nil {
			return err
		}

		// Step 4.1: 向量检索
		retrievedTextChunks, err = a.KnowledgeBaseRepo.SearchAgentRelatedChunks(
			c, int64(r.KnowledgeBaseID), embedding[0])
		if err != nil {
			zap.L().Error("向量检索失败", zap.Error(err))
			return fmt.Errorf("向量检索失败: %w", err)
		}

		// Step 4.2: 重排序
		rerankModel, err := a.AIModelRepo.GetRerankModel(c)
		if err != nil {
			zap.L().Error("获取重排序模型失败", zap.Error(err))
			return fmt.Errorf("获取重排序模型失败: %w", err)
		}
		rerankProvider, err := a.AIProviderRepo.GetAIProviderByIDWithCache(c, rerankModel.ProviderID)
		if err != nil {
			zap.L().Error("获取重排序模型提供者失败", zap.Error(err))
			return fmt.Errorf("获取重排序模型提供者失败: %w", err)
		}
		documents := make([]string, 0, len(retrievedTextChunks))
		for _, chunk := range retrievedTextChunks {
			documents = append(documents, chunk.Content)
		}
		rerankTexts, err := a.AIModelClient.Rerank(rerankProvider.APIKey, ai.RerankRequest{
			Model:     rerankModel.Model,
			Documents: documents,
			Query:     latestQuestion,
		}, 3)
		if err != nil {
			zap.L().Error("重排序失败", zap.Error(err))
			return fmt.Errorf("重排序失败: %w", err)
		}

		// Step 4.3 : 构造 prompt
		if len(documents) > 0 {
			contextText := strings.Join(rerankTexts, "\n\n")
			fullMessages = append(fullMessages, ai.ChatMessage{
				Role: "system",
				Content: fmt.Sprintf(`以下是与用户问题相关的背景资料。
						请仅根据这些资料回答问题，如果资料中没有提及，请直接回复：
						“很抱歉，我无法在现有知识中找到相关答案。”：%s`, contextText),
			})
		} else {
			fullMessages = append(fullMessages, ai.ChatMessage{
				Role:    "system",
				Content: `注意：当前没有任何与用户问题相关的背景资料。`,
			})
		}
	}

	fullMessages = append(fullMessages, r.Messages...)

	// Step 5: 在新对话中返回对话ID
	if r.ConversationID == 0 {
		var newConversationResp response.Conversation
		_ = copier.Copy(&newConversationResp, conversation)
		marshalString, _ := sonic.MarshalString(&newConversationResp)
		_, _ = fmt.Fprintf(c.Writer, "data: "+marshalString+"\n\n")
		c.Writer.(http.Flusher).Flush()
	}
	// Step 6: 调用大模型 (流式)
	AIResponse, err := a.AIModelClient.ChatStreamWithWriter(
		c.Request.Context(), c.Writer,
		provider.BaseURL+modelInfo.ApiPath,
		provider.APIKey,
		ai.DefaultStreamChatRequest(modelInfo.Model, fullMessages),
	)
	if err != nil {
		zap.L().Error("AI模型聊天失败", zap.Error(err))
		return
	}

	// Step 7: 异步保存会话
	// todo待测试
	go func(ctx context.Context, convID int64, question, answer string) {
		// 不要直接传 gin.Context
		prompt := &entity.Message{
			Role:           "user",
			Content:        question,
			ConversationID: convID,
		}
		assistantMessage := &entity.Message{
			Role:           "assistant",
			Content:        answer,
			ConversationID: convID,
		}

		if len(documents) > 0 {
			assistantMessage.KnowledgeBaseID = (*int64)(&r.KnowledgeBaseID)
			var fileChunkIDs []int64
			for _, chunk := range retrievedTextChunks {
				fileChunkIDs = append(fileChunkIDs, chunk.ID)
			}
			// 注意：FileChunkIDs 可能需要单独存表
			assistantMessage.FileChunkIDs = fileChunkIDs
		}

		if err := a.ConversationRepo.CreateMessage(ctx, prompt); err != nil {
			zap.L().Error("保存用户提问消息失败", zap.Error(err))
			return
		}
		if err := a.ConversationRepo.CreateMessage(ctx, assistantMessage); err != nil {
			zap.L().Error("保存AI回答消息失败", zap.Error(err))
			return
		}
	}(c.Copy(), conversation.ID, latestQuestion, AIResponse)

	return
}

func (a *AIService) GetHistoryConversation(c context.Context) (res []*response.Conversation, err error) {
	userID := authutils.GetUserID(c)
	historyConversations, err := a.ConversationRepo.GetHistoryConversation(c, userID)
	res = make([]*response.Conversation, 0, len(historyConversations))
	for _, conversation := range historyConversations {
		var conversationRes response.Conversation
		_ = copier.Copy(&conversationRes, conversation)
		res = append(res, &conversationRes)
	}
	return
}

func (a *AIService) GetMessageByConversationID(c context.Context, id int64) (res []*response.Message, err error) {
	messages, err := a.ConversationRepo.GetMessageByConversationID(c, id)
	res = make([]*response.Message, 0, len(messages))
	for _, message := range messages {
		var messageRes response.Message
		_ = copier.Copy(&messageRes, message)
		res = append(res, &messageRes)
	}
	return
}

//func (a *AIService) UploadAndVectorizeDocument(
//	c context.Context,
//	file multipart.File,
//	filename string,
//	agentID int64,
//) error {
//	// Step 1: 上传文档到 OSS
//	documentURL, err := a.OssClient.UploadAgentDocument(c, filename, file)
//	if err != nil {
//		zap.L().Error("上传文档失败", zap.Error(err))
//		return fmt.Errorf("上传文档失败: %w", err)
//	}
//
//	// Step 2: 创建 AgentFile 记录
//	agentFile := &entity.AgentFile{
//		AgentID: agentID,
//		FileURL: documentURL,
//	}
//	if err = a.AIAgentRepo.CreateAgentFile(c, agentFile); err != nil {
//		zap.L().Error("保存 AgentFile 失败", zap.Error(err))
//		return fmt.Errorf("保存 AgentFile 失败: %w", err)
//	}
//
//	// Step 3: 文档分块
//	chunks, err := a.Slicer.GetChunksFromSlicer(documentURL)
//	if err != nil {
//		zap.L().Error("文档分块失败", zap.Error(err))
//		return fmt.Errorf("文档分块失败: %w", err)
//	}
//	if len(chunks) == 0 {
//		return fmt.Errorf("文档内容为空，无法切分")
//	}
//
//	// 使用事务包裹整个处理流程
//	err = a.GormTX.Transaction(c, func(ctx context.Context) error {
//		// Step 4: 保存分块内容
//		chunkList := make([]*entity.FileChunk, 0, len(chunks))
//		for i, chunk := range chunks {
//			chunkEntity := &entity.FileChunk{
//				FileID:     agentFile.ID,
//				ChunkIndex: i,
//				Content:    chunk,
//			}
//			if err = a.AIAgentRepo.CreateFileChunk(ctx, chunkEntity); err != nil {
//				zap.L().Error("创建文件块失败", zap.Int("index", i), zap.Error(err))
//				return fmt.Errorf("创建文件块失败(index %d): %w", i, err)
//			}
//			chunkList = append(chunkList, chunkEntity)
//		}
//
//		// Step 5: 分批获取 Embedding（每次最多 10 个 chunk）
//		batchSize := 10
//		allEmbeddings := make([][]float32, 0, len(chunkList))
//		for i := 0; i < len(chunks); i += batchSize {
//			end := i + batchSize
//			if end > len(chunks) {
//				end = len(chunks)
//			}
//			batchChunks := chunks[i:end]
//
//			// 调用千问 Embedding API（每次最多 10 个）
//			batchEmbeddings, err := a.getQianwenEmbeddings(ctx, batchChunks)
//			if err != nil {
//				return fmt.Errorf("获取 Embedding 失败(batch %d-%d): %w", i, end-1, err)
//			}
//			allEmbeddings = append(allEmbeddings, batchEmbeddings...)
//		}
//
//		// Step 6: 检查向量数量是否匹配
//		if len(allEmbeddings) != len(chunkList) {
//			zap.L().Error("向量数量与分块数量不一致",
//				zap.Int("chunks", len(chunkList)),
//				zap.Int("vectors", len(allEmbeddings)))
//			return fmt.Errorf("向量数量与分块数量不一致")
//		}
//
//		// Step 7: 保存向量
//		for i, chunk := range chunkList {
//			chunkVector := &entity.ChunkVector{
//				ChunkID:   chunk.ID,
//				Embedding: pgvector.NewVector(allEmbeddings[i]),
//			}
//			if err = a.AIAgentRepo.CreateChunkVector(ctx, chunkVector); err != nil {
//				zap.L().Error("保存向量失败", zap.Int64("chunkID", chunk.ID), zap.Error(err))
//				return fmt.Errorf("保存向量失败(chunkID %d): %w", chunk.ID, err)
//			}
//		}
//
//		// ✅ 输出日志（事务内）
//		zap.L().Info("文档上传与向量化完成",
//			zap.Int64("agentID", agentID),
//			zap.Int("chunkCount", len(chunks)),
//			zap.String("file", filename),
//		)
//		return nil
//	})
//
//	if err != nil {
//		// 事务已自动回滚，此处可补充额外日志
//		zap.L().Error("文档处理事务失败", zap.Error(err))
//		return err
//	}
//
//	return nil
//}

//func (a *AIService) CreateAgent(c *gin.Context, r *request.AIAgent) (err error) {
//	agent := &entity.AIAgent{}
//	_ = copier.Copy(agent, r)
//	if err = a.AIAgentRepo.Create(c, agent); err != nil {
//		zap.L().Error("创建AI模型失败", zap.Error(err))
//		return
//	}
//	return
//}
//
//func (a *AIService) GetAIAgentPage(c *gin.Context, q *query.AIAgent) (res *common.PaginationResp[*response.AIAgent], err error) {
//	agents, total, err := a.AIAgentRepo.GetAIAgentPage(c, q)
//	if err != nil {
//		zap.L().Error("获取AI模型分页失败", zap.Error(err))
//		return
//	}
//	agentRes := make([]*response.AIAgent, 0, len(agents))
//	for _, agent := range agents {
//		agentResp := &response.AIAgent{}
//		_ = copier.Copy(agentResp, agent)
//		agentRes = append(agentRes, agentResp)
//	}
//	res = common.BuildPageResp[*response.AIAgent](agentRes, total, q.PaginationReq)
//	return
//}
//
//func (a *AIService) GetSimpleAgentList(c context.Context) (res []*response.SimpleAIAgent, err error) {
//	var agentList []*entity.AIAgent
//	agentList, err = a.AIAgentRepo.GetSimpleAgentList(c)
//	if err != nil {
//		zap.L().Error("获取智能体列表失败", zap.Error(err))
//		return
//	}
//	res = make([]*response.SimpleAIAgent, 0, len(agentList))
//	for _, agent := range agentList {
//		agentResp := &response.SimpleAIAgent{}
//		_ = copier.Copy(agentResp, agent)
//		res = append(res, agentResp)
//	}
//	return
//}

//func (a *AIService) AgentChatCompletion(c *gin.Context, r *request.ChatCompletion) (err error) {
//	// Step 1: 获取 Agent 信息
//	agentInfo, err := a.AIAgentRepo.GetAIAgentByIDWithCache(c, int64(r.ID))
//	if err != nil {
//		zap.L().Error("获取智能体失败", zap.Error(err))
//		return fmt.Errorf("获取智能体失败: %w", err)
//	}
//
//	// Step 2: 获取用户最新提问（最后一条 user 消息）
//	var latestQuestion string
//	for i := len(r.Messages) - 1; i >= 0; i-- {
//		if r.Messages[i].Role == "user" {
//			latestQuestion = r.Messages[i].Content
//			break
//		}
//	}
//	if latestQuestion == "" {
//		return fmt.Errorf("未找到用户提问")
//	}
//
//	// Step 3: 向量化提问
//	embedding, err := a.getQianwenEmbeddings(c, []string{latestQuestion})
//	if err != nil {
//		return err
//	}
//
//	// Step 4: 基于 embedding 搜索相关内容（假设你有一个向量搜索接口）
//	retrievedTexts, err :=
//		a.KnowledgeBaseRepo.SearchAgentRelatedChunks(c, r.KnowledgeBaseID, embedding[0])
//	if err != nil {
//		zap.L().Error("向量检索失败", zap.Error(err))
//		return fmt.Errorf("向量检索失败: %w", err)
//	}
//
//	//Step 5: 使用重排序模型进一步优化
//	rerankModel, err := a.AIModelRepo.GetRerankModel(c)
//	if err != nil {
//		zap.L().Error("获取重排序模型失败", zap.Error(err))
//		return
//	}
//	rerankProvider, err := a.AIProviderRepo.GetAIProviderByIDWithCache(c, rerankModel.ProviderID)
//	if err != nil {
//		zap.L().Error("获取重排序模型提供者失败", zap.Error(err))
//		return
//	}
//	rerankTexts, err := a.AIModelClient.Rerank(rerankProvider.APIKey, ai.RerankRequest{
//		Model:     rerankModel.Model,
//		Documents: retrievedTexts,
//		Query:     latestQuestion,
//	}, 3)
//	if err != nil {
//		zap.L().Error("重排序失败", zap.Error(err))
//		return
//	}
//	zap.L().Info("向量并重排搜索结果", zap.String("query", latestQuestion), zap.Any("texts", rerankTexts))
//
//	// Step 6: 构造新的对话上下文（system + retrieved + user）
//	var fullMessages []ai.ChatMessage
//
//	// 添加 agent 的默认 system prompt（如有）
//	if agentInfo.SystemPrompt != "" {
//		fullMessages = append(fullMessages, ai.ChatMessage{
//			Role:    "system",
//			Content: agentInfo.SystemPrompt,
//		})
//	}
//
//	// 将检索到的文本合并为一段 context，并构造辅助 system 提示
//	if len(retrievedTexts) > 0 {
//		contextText := strings.Join(rerankTexts, "\n\n")
//		fullMessages = append(fullMessages, ai.ChatMessage{
//			Role: "system",
//			Content: fmt.Sprintf(`以下是与用户问题相关的背景资料。请仅根据这些资料回答问题，
//				如果资料中没有提及，请直接回复：“很抱歉，我无法在现有知识中找到相关答案。”：%s`, contextText),
//		})
//	} else {
//		fullMessages = append(fullMessages, ai.ChatMessage{
//			Role:    "system",
//			Content: `注意：当前没有任何与用户问题相关的背景资料。如果资料中没有提及，请直接回复：“很抱歉，我无法在现有知识中找到相关答案。”`,
//		})
//	}
//
//	// 拼接用户原始对话
//	fullMessages = append(fullMessages, r.Messages...)
//
//	// Step 7: 调用模型回答（根据 agentInfo.ModelID 决定使用哪个模型）
//	return a.ChatCompletion(c, &request.ChatCompletion{
//		ID:       common.LongStringID(agentInfo.ModelID),
//		Messages: fullMessages,
//	})
//}

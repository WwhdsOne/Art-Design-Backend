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
	"Art-Design-Backend/pkg/constant/prompt"
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

const multiModelID = 61331207874412809

type AIService struct {
	AIModelClient     *ai.AIModelClient             // AI客户端
	AIModelRepo       *repository.AIModelRepo       // 模型Repo
	AIProviderRepo    *repository.AIProviderRepo    // 模型供应商Repo
	KnowledgeBaseRepo *repository.KnowledgeBaseRepo // 知识库Repo
	ConversationRepo  *repository.ConversationRepo  // 会话Repo
	OssClient         *aliyun.OssClient             // 阿里云OSS
	GormTX            *db.GormTransactionManager    // 事务
}

// 获取嵌入向量
func (a *AIService) getQianwenEmbeddings(c context.Context, chunks []string) ([][]float32, error) {
	// qwen模型供应商ID
	// todo后续写到配置文件
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
		zap.L().Error("AI 模型已存在", zap.Error(err))
		return
	}
	if err = a.AIModelRepo.Create(c, &aiModel); err != nil {
		zap.L().Error("AI 模型创建失败", zap.Error(err))
		return
	}
	go func() {
		if redisErr := a.AIModelRepo.InvalidSimpleModelList(); redisErr != nil {
			zap.L().Error("模型信息简易列表失效失败", zap.Error(redisErr))
		}
	}()
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

	zap.L().Info("对话请求信息:", zap.Int64("conversation_id", int64(r.ConversationID)))

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
							Content: prompt.TitleSummaryPrompt,
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
	} else {
		conversation.ID = int64(r.ConversationID)
	}

	var textContext string
	var imageContext string
	// Step 4: 构建知识库回答
	var fullMessages []ai.ChatMessage
	documents := make([]string, 0)
	var retrievedTextChunks []*entity.FileChunk
	var embedding [][]float32

	// 4.1 向量检索
	if r.KnowledgeBaseID != 0 {
		embedding, err = a.getQianwenEmbeddings(c, []string{latestQuestion})
		if err != nil {
			return err
		}
		retrievedTextChunks, err = a.KnowledgeBaseRepo.SearchAgentRelatedChunks(c, int64(r.KnowledgeBaseID), embedding[0])
		if err != nil {
			zap.L().Error("向量检索失败", zap.Error(err))
			return fmt.Errorf("向量检索失败: %w", err)
		}

		// 重排序
		var rerankModel *entity.AIModel
		rerankModel, err = a.AIModelRepo.GetRerankModel(c)
		if err != nil {
			zap.L().Error("获取重排序模型失败", zap.Error(err))
			return fmt.Errorf("获取重排序模型失败: %w", err)
		}
		var rerankProvider *entity.AIProvider
		rerankProvider, err = a.AIProviderRepo.GetAIProviderByIDWithCache(c, rerankModel.ProviderID)
		if err != nil {
			zap.L().Error("获取重排序模型提供者失败", zap.Error(err))
			return fmt.Errorf("获取重排序模型提供者失败: %w", err)
		}

		for _, chunk := range retrievedTextChunks {
			documents = append(documents, chunk.Content)
		}
		rerankTexts := make([]string, 0)
		rerankTexts, err = a.AIModelClient.Rerank(rerankProvider.APIKey, ai.RerankRequest{
			Model:     rerankModel.Model,
			Documents: documents,
			Query:     latestQuestion,
		}, 3)
		if err != nil {
			zap.L().Error("重排序失败", zap.Error(err))
			return fmt.Errorf("重排序失败: %w", err)
		}

		if len(rerankTexts) > 0 {
			textContext = strings.Join(rerankTexts, "\n\n")
		}
	}

	// 4.2 图片理解
	if len(r.Files) > 0 {
		var multiModel *entity.AIModel
		multiModel, err = a.AIModelRepo.GetAIModelByID(c, multiModelID)
		if err != nil {
			zap.L().Error("获取多模态模型失败", zap.Error(err))
			return
		}
		var multiModelProvider *entity.AIProvider
		multiModelProvider, err = a.AIProviderRepo.GetAIProviderByIDWithCache(c, multiModel.ProviderID)
		if err != nil {
			zap.L().Error("获取多模态模型供应商失败", zap.Error(err))
			return
		}

		var multiModeMessages []ai.MultiModeChatMessage
		multiModeMessages = append(multiModeMessages, ai.MultiModeChatMessage{
			Role: "system",
			Content: []ai.MultiModeChatContent{
				{Type: "text", Text: prompt.ImageSummaryPrompt},
			},
		})
		for _, file := range r.Files {
			multiModeMessages = append(multiModeMessages, ai.MultiModeChatMessage{
				Role: "user",
				Content: []ai.MultiModeChatContent{
					{Type: "image_url", ImageUrl: file},
				},
			})
		}
		var imageSummaryResponse []byte
		imageSummaryResponse, err = a.AIModelClient.MultiModeChatRequest(
			c.Request.Context(),
			multiModelProvider.BaseURL+multiModel.ApiPath,
			multiModelProvider.APIKey,
			ai.DefaultMultiModeChatRequest(multiModel.Model, multiModeMessages),
		)
		if err != nil {
			zap.L().Error("图片理解失败", zap.Error(err))
		} else {
			var imageSummary ai.ChatCompletionResponse
			_ = sonic.Unmarshal(imageSummaryResponse, &imageSummary)
			imageContext = imageSummary.FirstText()
			zap.L().Info("图片理解成功", zap.Int64("conversation_id", conversation.ID))
		}
	}

	// 4.3 构造 system 消息
	if textContext != "" || imageContext != "" {
		fullMessages = append(fullMessages, ai.ChatMessage{
			Role: "system",
			Content: fmt.Sprintf(`以下是与用户问题相关的背景资料，请严格按照规则回答：
				1. 文本知识（来自知识库和向量检索）：
				%s
				
				2. 图片理解知识（来自用户上传的图片，多模态分析结果）：
				%s
				
				规则：
				- 请仅根据上述资料回答问题。
				- 如果用户提问的内容在以上资料中都没有提及，请直接回答：
				  “很抱歉，我无法在现有知识中找到相关答案。”
				- 不要自己推测或者添加额外信息。
				- 尽量用中文简洁自然地回答。
				`, textContext, imageContext),
		})
		fmt.Printf(`以下是与用户问题相关的背景资料，请严格按照规则回答：
				1. 文本知识（来自知识库和向量检索）：
				%s
				
				2. 图片理解知识（来自用户上传的图片，多模态分析结果）：
				%s
				
				规则：
				- 请仅根据上述资料回答问题。
				- 如果用户提问的内容在以上资料中都没有提及，请直接回答：
				  “很抱歉，我无法在现有知识中找到相关答案。”
				- 不要自己推测或者添加额外信息。
				- 尽量用中文简洁自然地回答。
				`, textContext, imageContext)
	} else {
		fullMessages = append(fullMessages, ai.ChatMessage{
			Role:    "system",
			Content: `注意：当前没有任何与用户问题相关的背景资料。`,
		})
	}

	// Step 5: 追加用户原始消息
	fullMessages = append(fullMessages, r.Messages...)

	// Step 6: 在新对话中返回对话ID
	if r.ConversationID == 0 {
		var newConversationResp response.Conversation
		_ = copier.Copy(&newConversationResp, conversation)
		marshalString, _ := sonic.MarshalString(&newConversationResp)
		_, _ = fmt.Fprintf(c.Writer, "data: "+marshalString+"\n\n")
		c.Writer.(http.Flusher).Flush()
	}
	// Step 7: 调用大模型 (流式)
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

	// Step 8: 异步保存会话
	// todo待测试
	go func(ctx context.Context, convID int64, question, answer string) {
		// 不要直接传 gin.Context
		userMessage := &entity.Message{
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
			assistantMessage.FileChunkIDs = fileChunkIDs
		}

		if err := a.ConversationRepo.CreateMessage(ctx, userMessage); err != nil {
			zap.L().Error("保存用户提问消息失败", zap.Error(err))
			return
		}
		if err := a.ConversationRepo.CreateMessage(ctx, assistantMessage); err != nil {
			zap.L().Error("保存AI回答消息失败", zap.Error(err))
			return
		}
		zap.L().Info("保存会话成功",
			zap.Int64("conversation_id", conversation.ID),
			zap.Int64("user_id", authutils.GetUserID(c)),
			zap.String("question", latestQuestion),
			zap.String("answer", AIResponse),
		)
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

func (a *AIService) UploadChatMessageImage(c *gin.Context, filename string, src multipart.File) (url string, err error) {
	url, err = a.OssClient.UploadChatMessageImage(c, filename, src)
	if err != nil {
		zap.L().Error("上传对话图片失败", zap.Error(err))
		return
	}
	return
}

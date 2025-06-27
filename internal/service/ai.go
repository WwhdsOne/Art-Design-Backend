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
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"mime/multipart"
)

type AIService struct {
	AIModelRepo    *repository.AIModelRepo    // 模型Repo
	AIModelClient  *ai.AIModelClient          // 聊天
	AIProviderRepo *repository.AIProviderRepo // 模型供应商Repo
	OssClient      *aliyun.OssClient          // 阿里云OSS
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
func (a *AIService) GetSimpleModelList(c context.Context) (res []*response.SimpleAIModel, err error) {
	res, err = a.AIModelRepo.GetSimpleModelList(c)
	if err != nil {
		zap.L().Error("GetSimpleModelList 失败", zap.Error(err))
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

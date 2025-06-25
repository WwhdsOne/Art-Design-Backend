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
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type AIService struct {
	AIModelRepo    *repository.AIModelRepo    // 模型Repo
	AIModelClient  *ai.AIModelClient          // 聊天
	AIProviderRepo *repository.AIProviderRepo // 模型供应商Repo
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
		zap.L().Error("AI模型创建失败", zap.Error(err))
		return
	}
	return
}
func (a *AIService) GetSimpleModelList(c context.Context) (res []*response.SimpleAIModel, err error) {
	res, err = a.AIModelRepo.GetSimpleModelList(c)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return
}

func (a *AIService) GetAIModelPage(c context.Context, q *query.AIModel) (res *common.PaginationResp[*response.AIModel], err error) {
	var aiModel []*entity.AIModel
	var total int64
	if aiModel, total, err = a.AIModelRepo.GetAIModelPage(c, q); err != nil {
		zap.L().Error("AI模型分页查询失败", zap.Error(err))
		return
	}
	aiModelRes := make([]*response.AIModel, 0, len(aiModel))
	for _, model := range aiModel {
		var aiModelResp response.AIModel
		_ = copier.Copy(&aiModelResp, model)
		aiModelRes = append(aiModelRes, &aiModelResp)
	}
	res = common.BuildPageResp[*response.AIModel](aiModelRes, total, q.PaginationReq)
	return
}

// ChatCompletion
// todo
// 1. 对话管理
// 2. 模型信息缓存 ✅
// 3. 标签抽取
// 4. 价格计算
//func (a *AIService) ChatCompletion(c *gin.Context, r *request.ChatCompletion) (err error) {
//	modelID := int64(r.ID)
//	modelInfo, err := a.AIModelRepo.GetAIModelByID(c, modelID)
//	if err != nil {
//		zap.L().Error("获取AI模型失败", zap.Int64("model_id", modelID), zap.Error(err))
//		return
//	}
//
//	reqData := ai.DefaultStreamChatRequest(modelInfo.Model, r.Messages)
//
//	// 直接调用新版 ChatStreamWithWriter
//	err = a.AIModelClient.ChatStreamWithWriter(c.Request.Context(), c.Writer, modelInfo.BaseURL, modelInfo.APIKey, reqData)
//	if err != nil {
//		zap.L().Error("AI模型聊天失败", zap.Error(err))
//		return err
//	}
//
//	return nil
//}

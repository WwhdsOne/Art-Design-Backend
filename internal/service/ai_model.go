package service

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/internal/repository/db"
	"Art-Design-Backend/pkg/ai"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type AIModelService struct {
	AIModelRepo   *repository.AIModelRepo    // 模型Repo
	AIModelClient *ai.AIModelClient          // 聊天
	GormTX        *db.GormTransactionManager // 事务
}

func (a *AIModelService) CreateAIModel(c context.Context, r *request.AIModel) (err error) {
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
func (a *AIModelService) GetSimpleModelList(c context.Context) (res []*response.SimpleAIModel, err error) {
	res, err = a.AIModelRepo.GetSimpleModelList(c)
	if err != nil {
		zap.L().Error("获取AI模型简洁信息列表失败", zap.Error(err))
	}
	return
}

func (a *AIModelService) GetAIModelPage(c context.Context, q *query.AIModel) (res *base.PaginationResp[*response.AIModel], err error) {
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
	res = base.BuildPageResp[*response.AIModel](aiModelRes, total, q.PaginationReq)
	return
}

// ChatCompletion
// todo
// 1. 对话管理
// 2. 模型信息缓存 ✅
// 3. 标签抽取
// 4. 价格计算
func (a *AIModelService) ChatCompletion(c *gin.Context, r *request.ChatCompletion) (err error) {
	modelID := int64(r.ID)
	modelInfo, err := a.AIModelRepo.GetAIModelByID(c, modelID)
	if err != nil {
		zap.L().Error("获取AI模型失败", zap.Int64("model_id", modelID), zap.Error(err))
		return
	}
	if err = a.AIModelClient.ChatStream(c,
		modelInfo.BaseURL,
		modelInfo.APIKey,
		ai.DefaultStreamChatRequest(modelInfo.Model, r.Messages)); err != nil {
		zap.L().Error("AI模型聊天失败", zap.Error(err))
		return
	}
	return
}

package service

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/ai"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type AIModelService struct {
	AIModelRepo   *repository.AIModelRepository      // AI模型
	GormTX        *repository.GormTransactionManager // 事务
	AIModelClient *ai.AIModelClient                  // 聊天
	Redis         *redisx.RedisWrapper               // redis
}

func (a *AIModelService) CreateAIModel(c context.Context, r *request.AIModel) (err error) {
	var aiModel entity.AIModel
	if err = copier.Copy(&aiModel, &r); err != nil {
		zap.L().Error("AI模型属性复制失败", zap.Error(err))
		return
	}
	if err = a.AIModelRepo.CheckAIDuplicate(&aiModel); err != nil {
		zap.L().Error("AI模型已存在", zap.Error(err))
		return
	}
	if err = a.AIModelRepo.Create(c, &aiModel); err != nil {
		zap.L().Error("AI模型创建失败", zap.Error(err))
		return
	}
	return
}

// todo 后续添加缓存
func (a *AIModelService) GetSimpleModelList(c context.Context) (res []*response.SimpleAIModel, err error) {
	var aiModel []*entity.AIModel
	if aiModel, err = a.AIModelRepo.GetSimpleModelList(c); err != nil {
		zap.L().Error("获取AI模型列表失败", zap.Error(err))
		return
	}
	res = make([]*response.SimpleAIModel, 0, len(aiModel))
	for _, model := range aiModel {
		var aiModelResp response.SimpleAIModel
		if err = copier.Copy(&aiModelResp, model); err != nil {
			zap.L().Error("AI模型属性复制失败", zap.Error(err))
			return
		}
		res = append(res, &aiModelResp)
	}
	return
}

func (a *AIModelService) GetAIModelPage(c context.Context, q *query.AIModel) (res base.PaginationResp[*response.AIModel], err error) {
	var aiModel []*entity.AIModel
	var total int64
	if aiModel, total, err = a.AIModelRepo.GetAIModelPage(c, q); err != nil {
		zap.L().Error("AI模型分页查询失败", zap.Error(err))
		return
	}
	aiModelRes := make([]*response.AIModel, 0, len(aiModel))
	for _, model := range aiModel {
		var aiModelResp response.AIModel
		if err = copier.Copy(&aiModelResp, model); err != nil {
			zap.L().Error("AI模型属性复制失败", zap.Error(err))
			return
		}
		aiModelRes = append(aiModelRes, &aiModelResp)
	}
	res = *base.BuildPageResp[*response.AIModel](aiModelRes, total, q.PaginationReq)
	return
}

// todo 1.对话管理 2.模型信息缓存 3. 标签抽取 4. 价格计算
func (a *AIModelService) ChatCompletion(c *gin.Context, r *request.ChatCompletion) (err error) {
	res, err := a.AIModelRepo.GetAIModelByID(c, int64(r.ID))
	if err != nil {
		zap.L().Error("获取AI模型失败", zap.Error(err))
		return
	}
	if err = a.AIModelClient.ChatStream(c, res.BaseURL, res.APIKey, ai.DefaultStreamChatRequest(res.Model, r.Messages)); err != nil {
		zap.L().Error("AI模型聊天失败", zap.Error(err))
		return
	}
	return
}

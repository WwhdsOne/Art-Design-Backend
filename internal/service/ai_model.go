package service

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type AIModelService struct {
	AIModelRepo *repository.AIModelRepository      // AI模型
	GormTX      *repository.GormTransactionManager // 事务
	Redis       *redisx.RedisWrapper               // redis
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

func (a *AIModelService) GetAIModelPage(c *gin.Context, q *query.AIModel) (res base.PaginationResp[*response.AIModel], err error) {
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

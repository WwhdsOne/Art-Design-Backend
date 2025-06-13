package service

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"Art-Design-Backend/pkg/ai"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type AIModelService struct {
	AIModelRepo   *db.AIModelRepository      // AI模型
	AIModelCache  *cache.AIModelCache        // 缓存
	GormTX        *db.GormTransactionManager // 事务
	AIModelClient *ai.AIModelClient          // 聊天
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
func (a *AIModelService) GetSimpleModelList(c context.Context) (res []*response.SimpleAIModel, err error) {
	res, err = a.AIModelCache.GetSimpleModelList()
	if err == nil {
		return
	}
	zap.L().Warn("缓存获取失败或为空，尝试访问数据库", zap.Error(err))

	var aiModels []*entity.AIModel
	if aiModels, err = a.AIModelRepo.GetSimpleModelList(c); err != nil {
		zap.L().Error("数据库查询失败", zap.Error(err))
		return nil, err
	}

	res = make([]*response.SimpleAIModel, 0, len(aiModels))
	for _, model := range aiModels {
		var aiModelResp response.SimpleAIModel
		if err = copier.Copy(&aiModelResp, model); err != nil {
			zap.L().Error("AI模型属性复制失败，跳过该条", zap.Error(err))
			continue
		}
		res = append(res, &aiModelResp)
	}

	// 异步回写缓存（避免闭包共享 err）
	go func(models []*response.SimpleAIModel) {
		if err := a.AIModelCache.SetSimpleModelList(models); err != nil {
			zap.L().Error("设置AI模型列表缓存失败", zap.Error(err))
		}
	}(res)

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

// ChatCompletion
// todo
// 1. 对话管理
// 2. 模型信息缓存 ✅
// 3. 标签抽取
// 4. 价格计算
func (a *AIModelService) ChatCompletion(c *gin.Context, r *request.ChatCompletion) (err error) {
	var modelInfo *entity.AIModel
	modelID := int64(r.ID)
	modelInfo, err = a.AIModelCache.GetModelInfo(modelID)
	if err != nil {
		// 缓存中不存在
		zap.L().Warn("AI模型缓存获取失败", zap.Int64("model_id", modelID), zap.Error(err))
		modelInfo, err = a.AIModelRepo.GetAIModelByID(c, modelID)
		if err != nil {
			zap.L().Error("获取AI模型失败", zap.Int64("model_id", modelID), zap.Error(err))
			return
		}
		go func(m *entity.AIModel) {
			if err := a.AIModelCache.SetModelInfo(m); err != nil {
				zap.L().Error("设置AI模型缓存失败", zap.Int64("model_id", m.ID), zap.Error(err))
			}
		}(modelInfo)
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

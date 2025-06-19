package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AIModelRepo struct {
	aiModelDB    *db.AIModelDB
	aiModelCache *cache.AIModelCache
}

func NewAIModelRepo(
	aiModelDB *db.AIModelDB,
	aiModelCache *cache.AIModelCache,
) *AIModelRepo {
	return &AIModelRepo{
		aiModelDB:    aiModelDB,
		aiModelCache: aiModelCache,
	}
}

func (a *AIModelRepo) CheckAIDuplicate(c context.Context, model *entity.AIModel) (err error) {
	return a.aiModelDB.CheckAIDuplicate(c, model)
}

func (a *AIModelRepo) Create(c context.Context, e *entity.AIModel) (err error) {
	return a.aiModelDB.Create(c, e)
}

func (a *AIModelRepo) GetAIModelByID(c context.Context, id int64) (res *entity.AIModel, err error) {
	res, err = a.aiModelCache.GetModelInfo(id)
	if err != nil {
		zap.L().Warn("获取AI模型缓存失败", zap.Error(err))
	}
	res, err = a.aiModelDB.GetAIModelByID(c, id)
	if err != nil {
		return
	}
	go func() {
		_ = a.aiModelCache.SetModelInfo(res)
	}()
	return
}

func (a *AIModelRepo) GetAIModelPage(c context.Context, q *query.AIModel) (pageRes []*entity.AIModel, total int64, err error) {
	return a.aiModelDB.GetAIModelPage(c, q)
}

func (a *AIModelRepo) GetSimpleModelList(c context.Context) (res []*response.SimpleAIModel, err error) {
	res, err = a.aiModelCache.GetSimpleModelList()
	if err != nil && !errors.Is(err, redis.Nil) {
		zap.L().Warn("获取AI模型简洁信息列表缓存失败", zap.Error(err))
	}
	list, err := a.aiModelDB.GetSimpleModelList(c)
	if err != nil {
		return
	}
	// 类型转换
	res = make([]*response.SimpleAIModel, 0, len(list))
	for _, model := range list {
		var simpleModel response.SimpleAIModel
		_ = copier.Copy(&simpleModel, model)
		res = append(res, &simpleModel)
	}
	// 写入缓存
	go func(model []*response.SimpleAIModel) {
		_ = a.aiModelCache.SetSimpleModelList(model)
	}(res)
	return
}

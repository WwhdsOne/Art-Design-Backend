package repository

import (
	"Art-Design-Backend/internal/model/entity"
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
	*db.AIModelDB
	*cache.AIModelCache
}

func (a *AIModelRepo) GetAIModelByID(c context.Context, id int64) (res *entity.AIModel, err error) {
	res, err = a.AIModelCache.GetModelInfo(id)
	if err != nil {
		zap.L().Warn("获取AI模型缓存失败", zap.Error(err))
	}
	res, err = a.AIModelDB.GetAIModelByID(c, id)
	if err != nil {
		return
	}
	go func(res *entity.AIModel) {
		if err = a.AIModelCache.SetModelInfo(res); err != nil {
			zap.L().Warn("设置AI模型缓存失败", zap.Error(err))
		}
	}(res)
	return
}

func (a *AIModelRepo) GetSimpleModelList(c context.Context) (res []*response.SimpleAIModel, err error) {
	res, err = a.AIModelCache.GetSimpleModelList()
	if err != nil && !errors.Is(err, redis.Nil) {
		zap.L().Warn("获取AI模型简洁信息列表缓存失败", zap.Error(err))
	}
	list, err := a.AIModelDB.GetSimpleModelList(c)
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
		if err = a.AIModelCache.SetSimpleModelList(model); err != nil {
			zap.L().Warn("设置AI模型简洁信息列表缓存失败", zap.Error(err))
		}
	}(res)
	return
}

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
	if err == nil {
		// 缓存命中，直接返回
		return
	}
	if !errors.Is(err, redis.Nil) {
		// 缓存出错，但不是未命中，记录日志
		zap.L().Warn("获取AI模型缓存失败", zap.Error(err))
	}

	// 缓存未命中，查数据库
	res, err = a.AIModelDB.GetAIModelByID(c, id)
	if err != nil {
		return nil, err
	}

	// 异步回填缓存
	go func(model *entity.AIModel) {
		if model == nil {
			return
		}
		if cacheErr := a.AIModelCache.SetModelInfo(model); cacheErr != nil {
			zap.L().Warn("设置AI模型缓存失败", zap.Error(cacheErr))
		}
	}(res)

	return res, nil
}

func (a *AIModelRepo) GetSimpleChatModelList(c context.Context) (res []*response.SimpleAIModel, err error) {
	res, err = a.AIModelCache.GetSimpleModelList()
	if err == nil {
		// 命中缓存，返回结果
		return
	}
	if !errors.Is(err, redis.Nil) {
		// 如果不是缓存未命中，而是真正的错误，记录日志
		zap.L().Warn("获取AI模型简洁信息列表缓存失败", zap.Error(err))
	}

	// 缓存未命中，继续执行后续逻辑
	list, err := a.AIModelDB.GetSimpleChatModelList(c)
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

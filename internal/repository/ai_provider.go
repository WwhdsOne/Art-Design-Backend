package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AIProviderRepo struct {
	*db.AIProviderDB
	*cache.AIProviderCache
}

func (a *AIProviderRepo) GetAIProviderByIDWithCache(c context.Context, id int64) (res *entity.AIProvider, err error) {
	res, err = a.GetProviderCacheByID(id)
	if err == nil {
		// 缓存命中，直接返回
		return
	}
	if !errors.Is(err, redis.Nil) {
		// 缓存出错，但不是未命中，记录日志
		zap.L().Warn("根据ID查询供应商缓存失败", zap.Error(err))
	} else {
		zap.L().Warn("未命中供应商缓存")
	}

	// 缓存未命中，继续查数据库
	res, err = a.GetProviderByID(c, id)
	if err != nil {
		return
	}
	go func(provider *entity.AIProvider) {
		if err := a.SetProviderCacheByID(provider); err != nil {
			zap.L().Warn("设置供应商缓存失败", zap.Error(err))
		}
	}(res)
	return
}

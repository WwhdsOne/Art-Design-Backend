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
	if err != nil && !errors.Is(err, redis.Nil) {
		zap.L().Warn("根据ID查询供应商缓存失败", zap.Error(err))
	} else {
		return
	}
	provider, err := a.GetProviderByID(c, id)
	if err != nil {
		return
	}
	go func(provider *entity.AIProvider) {
		if err := a.SetProviderCacheByID(provider); err != nil {
			zap.L().Warn("设置供应商缓存失败", zap.Error(err))
		}
	}(provider)
	return
}

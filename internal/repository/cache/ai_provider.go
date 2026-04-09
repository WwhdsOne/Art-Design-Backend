package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/cachex"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/redisx"
	"strconv"
)

type AIProviderCache struct {
	redis *redisx.RedisWrapper
}

func NewAIProviderCache(redis *redisx.RedisWrapper) *AIProviderCache {
	return &AIProviderCache{
		redis: redis,
	}
}

func (a *AIProviderCache) GetProviderCacheByID(providerID int64) (*entity.AIProvider, error) {
	return cachex.GetOrNil[entity.AIProvider](a.redis, rediskey.AIProviderInfo+strconv.FormatInt(providerID, 10))
}

func (a *AIProviderCache) SetProviderCacheByID(provider *entity.AIProvider) error {
	return cachex.Set(a.redis, rediskey.AIProviderInfo+strconv.FormatInt(provider.ID, 10), provider, rediskey.AIProviderInfoTTL)
}

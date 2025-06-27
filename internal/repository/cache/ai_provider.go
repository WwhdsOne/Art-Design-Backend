package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/redisx"
	"github.com/bytedance/sonic"
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

func (a *AIProviderCache) GetProviderCacheByID(providerID int64) (res *entity.AIProvider, err error) {
	val, err := a.redis.Get(rediskey.AIProviderInfo + strconv.FormatInt(providerID, 10))
	if err != nil {
		return
	}
	_ = sonic.Unmarshal([]byte(val), &res)
	return
}

func (a *AIProviderCache) SetProviderCacheByID(provider *entity.AIProvider) (err error) {
	val, err := sonic.Marshal(provider)
	if err != nil {
		return
	}
	err = a.redis.Set(rediskey.AIProviderInfo+strconv.FormatInt(provider.ID, 10), string(val), rediskey.AIProviderInfoTTL)
	return
}

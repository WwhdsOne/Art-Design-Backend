package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository/cachex"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/redisx"
	"strconv"
)

type AIModelCache struct {
	redis *redisx.RedisWrapper
}

func NewAIModelCache(redis *redisx.RedisWrapper) *AIModelCache {
	return &AIModelCache{
		redis: redis,
	}
}

func (a *AIModelCache) GetSimpleModelList() ([]*response.SimpleAIModel, error) {
	return cachex.GetSlice[*response.SimpleAIModel](a.redis, rediskey.AIModelSimpleList)
}

func (a *AIModelCache) SetSimpleModelList(res []*response.SimpleAIModel) error {
	return cachex.Set(a.redis, rediskey.AIModelSimpleList, res, rediskey.AIModelSimpleListTTL)
}

func (a *AIModelCache) InvalidSimpleModelList() error {
	return a.redis.Del(rediskey.AIModelSimpleList)
}

func (a *AIModelCache) GetModelInfo(modelID int64) (*entity.AIModel, error) {
	return cachex.GetOrNil[entity.AIModel](a.redis, rediskey.AIModelInfo+strconv.FormatInt(modelID, 10))
}

func (a *AIModelCache) SetModelInfo(model *entity.AIModel) error {
	return cachex.Set(a.redis, rediskey.AIModelInfo+strconv.FormatInt(model.ID, 10), model, rediskey.AIModelInfoTTL)
}

func (a *AIModelCache) InvalidModelInfo(modelID int64) error {
	return a.redis.Del(rediskey.AIModelInfo + strconv.FormatInt(modelID, 10))
}

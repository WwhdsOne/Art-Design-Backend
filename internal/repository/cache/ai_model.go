package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/redisx"
	"strconv"

	"github.com/bytedance/sonic"
)

type AIModelCache struct {
	redis *redisx.RedisWrapper
}

func NewAIModelCache(redis *redisx.RedisWrapper) *AIModelCache {
	return &AIModelCache{
		redis: redis,
	}
}

func (a *AIModelCache) GetSimpleModelList() (res []*response.SimpleAIModel, err error) {
	val, err := a.redis.Get(rediskey.AIModelSimpleList)
	if err != nil {
		return
	}
	if err = sonic.Unmarshal([]byte(val), &res); err != nil {
		return
	}
	return
}

func (a *AIModelCache) SetSimpleModelList(res []*response.SimpleAIModel) (err error) {
	val, err := sonic.Marshal(res)
	if err != nil {
		return
	}
	err = a.redis.Set(rediskey.AIModelSimpleList, string(val), rediskey.AIModelSimpleListTTL)
	return
}

func (a *AIModelCache) InvalidSimpleModelList() error {
	return a.redis.Del(rediskey.AIModelSimpleList)
}

func (a *AIModelCache) GetModelInfo(modelID int64) (res *entity.AIModel, err error) {
	val, err := a.redis.Get(rediskey.AIModelInfo + strconv.FormatInt(modelID, 10))
	if err != nil {
		return
	}
	if err = sonic.Unmarshal([]byte(val), &res); err != nil {
		return
	}
	return
}

func (a *AIModelCache) SetModelInfo(model *entity.AIModel) (err error) {
	val, err := sonic.Marshal(model)
	if err != nil {
		return
	}
	err = a.redis.Set(rediskey.AIModelInfo+strconv.FormatInt(model.ID, 10), string(val), rediskey.AIModelInfoTTL)
	return
}

func (a *AIModelCache) InvalidModelInfo(modelID int64) (err error) {
	err = a.redis.Del(rediskey.AIModelInfo + strconv.FormatInt(modelID, 10))
	return
}

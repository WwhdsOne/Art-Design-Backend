package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/pkg/redisx"
	"github.com/bytedance/sonic"
	"strconv"
	"time"
)

const (
	AIModelSimpleList    = "AIMODEL:SIMPLE:LIST"
	AIModelSimpleListTTL = 86400 * time.Second
	AIModelInfo          = "AIMODEL:INFO:"
	AIModelInfoTTL       = 86400 * time.Second
)

type AIModelCache struct {
	Redis *redisx.RedisWrapper
}

func NewAIModelCache(redis *redisx.RedisWrapper) *AIModelCache {
	return &AIModelCache{
		Redis: redis,
	}
}

func (a *AIModelCache) GetSimpleModelList() (res []*response.SimpleAIModel, err error) {
	val, err := a.Redis.Get(AIModelSimpleList)
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
	err = a.Redis.Set(AIModelSimpleList, string(val), AIModelSimpleListTTL)
	return
}

func (a *AIModelCache) InvalidSimpleModelList() error {
	return a.Redis.Del(AIModelSimpleList)
}

func (a *AIModelCache) GetModelInfo(modelID int64) (res *entity.AIModel, err error) {
	val, err := a.Redis.Get(AIModelInfo + strconv.FormatInt(modelID, 10))
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
	err = a.Redis.Set(AIModelInfo+strconv.FormatInt(model.ID, 10), string(val), AIModelInfoTTL)
	return
}

func (a *AIModelCache) InvalidModelInfo(modelID int64) (err error) {
	err = a.Redis.Del(AIModelInfo + strconv.FormatInt(modelID, 10))
	return
}

package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/redisx"
	"github.com/bytedance/sonic"
	"strconv"
)

type AIAgentCache struct {
	redis *redisx.RedisWrapper
}

func NewAIAgentCache(redis *redisx.RedisWrapper) *AIAgentCache {
	return &AIAgentCache{
		redis: redis,
	}
}

func (a *AIAgentCache) GetAIAgentInfo(agentID int64) (res *entity.AIAgent, err error) {
	val, err := a.redis.Get(rediskey.AIAgentInfo + strconv.FormatInt(agentID, 10))
	if err != nil {
		return
	}
	if err = sonic.Unmarshal([]byte(val), &res); err != nil {
		return
	}
	return
}

func (a *AIAgentCache) SetAgentInfo(agent *entity.AIAgent) (err error) {
	val, err := sonic.Marshal(agent)
	if err != nil {
		return
	}
	err = a.redis.Set(rediskey.AIAgentInfo+strconv.FormatInt(agent.ID, 10), string(val), rediskey.AIAgentInfoTTL)
	return
}

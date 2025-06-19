package cache

import (
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/redisx"
	"strconv"
	"time"
)

type AuthCache struct {
	redis *redisx.RedisWrapper
}

func NewAuthCache(redis *redisx.RedisWrapper) *AuthCache {
	return &AuthCache{
		redis: redis,
	}
}

// GetTokenByUserID 获取用户token
func (a *AuthCache) GetTokenByUserID(userID int64) (token string, err error) {
	key := rediskey.SESSION + strconv.FormatInt(userID, 10)
	return a.redis.Get(key)
}

// DeleteOldSession 删除旧会话
func (a *AuthCache) DeleteOldSession(userID int64, token string) (err error) {
	sessionKey := rediskey.SESSION + strconv.FormatInt(userID, 10)
	loginKey := rediskey.LOGIN + token
	keys := []string{loginKey, sessionKey}
	return a.redis.PipelineDel(keys)
}

// SetNewSession 设置用户新的token
func (a *AuthCache) SetNewSession(userID int64, token string, ttl time.Duration) (err error) {
	loginKey := rediskey.LOGIN + token
	sessionKey := rediskey.SESSION + strconv.FormatInt(userID, 10)
	keyVals := [][2]string{
		{loginKey, strconv.FormatInt(userID, 10)},
		{sessionKey, token},
	}
	return a.redis.PipelineSet(keyVals, ttl)
}

// DeleteUserSession 登出
func (a *AuthCache) DeleteUserSession(userID int64, token string) (err error) {
	sessionKey := rediskey.SESSION + strconv.FormatInt(userID, 10)
	loginKey := rediskey.LOGIN + token
	delKeys := []string{sessionKey, loginKey}
	return a.redis.PipelineDel(delKeys)
}

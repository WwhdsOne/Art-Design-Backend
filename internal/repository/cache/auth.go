package cache

import (
	"Art-Design-Backend/pkg/redisx"
	"strconv"
	"time"
)

const (
	LOGIN   = "AUTH:LOGIN:"   // token -> userID
	SESSION = "AUTH:SESSION:" // userID -> token
)

type AuthCache struct {
	Redis *redisx.RedisWrapper
}

func NewAuthCache(redis *redisx.RedisWrapper) *AuthCache {
	return &AuthCache{
		Redis: redis,
	}
}

// GetTokenByUserID 获取用户token
func (a *AuthCache) GetTokenByUserID(userID int64) (token string, err error) {
	key := SESSION + strconv.FormatInt(userID, 10)
	return a.Redis.Get(key)
}

// DeleteOldSession 删除旧会话
func (a *AuthCache) DeleteOldSession(userID int64, token string) (err error) {
	sessionKey := SESSION + strconv.FormatInt(userID, 10)
	loginKey := LOGIN + token
	keys := []string{loginKey, sessionKey}
	return a.Redis.PipelineDel(keys)
}

// SetNewSession 设置用户新的token
func (a *AuthCache) SetNewSession(userID int64, token string, ttl time.Duration) (err error) {
	loginKey := LOGIN + token
	sessionKey := SESSION + strconv.FormatInt(userID, 10)
	keyVals := [][2]string{
		{loginKey, strconv.FormatInt(userID, 10)},
		{sessionKey, token},
	}
	return a.Redis.PipelineSet(keyVals, ttl)
}

// LogoutToken 登出
func (a *AuthCache) LogoutToken(userID int64, token string) (err error) {
	sessionKey := SESSION + strconv.FormatInt(userID, 10)
	loginKey := LOGIN + token
	delKeys := []string{sessionKey, loginKey}
	return a.Redis.PipelineDel(delKeys)
}

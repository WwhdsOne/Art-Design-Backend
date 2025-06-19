package repository

import (
	"Art-Design-Backend/internal/repository/cache"
	"go.uber.org/zap"
	"time"
)

type AuthRepo struct {
	authCache *cache.AuthCache
}

func NewAuthRepo(authCache *cache.AuthCache) *AuthRepo {
	return &AuthRepo{
		authCache: authCache,
	}
}

func (a *AuthRepo) GetTokenByUserID(userID int64) (string, error) {
	return a.authCache.GetTokenByUserID(userID)
}

func (a *AuthRepo) SetNewSession(userID int64, token string, ttl time.Duration) error {
	return a.authCache.SetNewSession(userID, token, ttl)
}

// LogoutByUserID 登出用户（先查 token 再删）
func (a *AuthRepo) LogoutByUserID(userID int64) (err error) {
	var token string
	token, err = a.authCache.GetTokenByUserID(userID)
	if err != nil {
		zap.L().Warn("获取用户 token 失败，可能已登出", zap.Int64("userID", userID), zap.Error(err))
	}
	return a.authCache.DeleteUserSession(userID, token)
}

func (a *AuthRepo) DeleteOldSession(userID int64, token string) (err error) {
	return a.authCache.DeleteOldSession(userID, token)
}

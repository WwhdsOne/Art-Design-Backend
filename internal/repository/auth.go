package repository

import (
	"Art-Design-Backend/internal/repository/cache"

	"go.uber.org/zap"
)

type AuthRepo struct {
	*cache.AuthCache
}

// LogoutByUserID 登出用户（先查 token 再删）
func (a *AuthRepo) LogoutByUserID(userID int64) (err error) {
	var token string
	token, err = a.AuthCache.GetTokenByUserID(userID)
	if err != nil {
		zap.L().Warn("获取用户 token 失败，可能已登出", zap.Int64("userID", userID), zap.Error(err))
	}
	return a.AuthCache.DeleteUserSession(userID, token)
}

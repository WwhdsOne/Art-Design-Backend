package middleware

import (
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/result"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (m *Middlewares) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Token
		token := authutils.GetToken(c)
		if token == "" {
			result.NoAuth("缺少 Token", c)
			c.Abort()
			return
		}

		// 2. 解析 Token
		claims, err := m.Jwt.ParseToken(token)
		if err != nil {
			handleTokenError(err, token, c)
			c.Abort()
			return
		}

		// 3. 校验 Redis 中的登录状态
		_, err = m.Redis.Get(rediskey.LOGIN + token)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				zap.L().Info("用户未登录", zap.String("token", token))
				result.NoAuth("当前未登录", c)
			} else {
				zap.L().Error("Redis 获取 Session 错误", zap.String("token", token), zap.Error(err))
				result.FailWithMessage("获取 Session 失败", c)
			}
			c.Abort()
			return
		}

		// 4. 认证通过，设置 claims 并继续请求
		c.Set("claims", claims)
		c.Next()
	}
}

func (m *Middlewares) WSAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Query("token")
		if token == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Abort()
			return
		}

		claims, err := m.Jwt.ParseToken(token)
		if err != nil {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Abort()
			return
		}

		_, err = m.Redis.Get(rediskey.LOGIN + token)
		if err != nil {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

// 辅助函数：处理 Token 错误
func handleTokenError(err error, token string, c *gin.Context) {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		result.NoAuth("授权已过期", c)
	case errors.Is(err, jwt.ErrTokenNotValidYet),
		errors.Is(err, jwt.ErrTokenMalformed),
		errors.Is(err, jwt.ErrTokenInvalid):
		result.NoAuth("Token 无效", c)
	default:
		zap.L().Error("Token 解析失败", zap.String("token", token), zap.Error(err))
		result.FailWithMessage("Token 解析失败", c)
	}
}

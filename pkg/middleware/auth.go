package middleware

import (
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/result"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (m *Middlewares) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := authutils.GetToken(c)
		if token == "" {
			result.NoAuth("缺少Token", c)
			c.Abort()
			return
		}

		// 校验登录状态（Redis中是否存在token对应会话）
		_, err := m.Redis.Get(rediskey.LOGIN + token)
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
		zap.L().Info("Auth Token 验证成功", zap.String("token", token))

		// 解析 Token
		claims, err := m.Jwt.ParseToken(token)
		if err != nil {
			switch {
			case errors.Is(err, jwt.TokenExpired):
				result.NoAuth("授权已过期", c)
			case errors.Is(err, jwt.TokenNotValidYet),
				errors.Is(err, jwt.TokenMalformed),
				errors.Is(err, jwt.TokenInvalid):
				result.NoAuth("token无效", c)
			default:
				zap.L().Error("Token 解析失败", zap.String("token", token), zap.Error(err))
				result.FailWithMessage("Token 解析失败", c)
			}
			c.Abort()
			return
		}

		// 判断是否需要刷新token（排除登出接口）
		if jwt.IsWithinRefreshWindow(claims) && c.FullPath() != "/api/auth/logout" {
			result.ShouldRefresh(c)
			c.Abort()
			return
		}

		// 设置解析后的 claims，继续后续处理
		c.Set("claims", claims)
		c.Next()
	}
}

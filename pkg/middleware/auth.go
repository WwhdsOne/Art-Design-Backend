package middleware

import (
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/result"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (m *Middlewares) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息
		// 这里前端需要把token存储到cookie或者本地localStorage中
		// 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := authutils.GetToken(c)
		id := m.Redis.Get(rediskey.LOGIN + token)
		if id != "" {
			zap.L().Info(fmt.Sprintf("Auth Token: %s", token))
		} else {
			zap.L().Error(fmt.Sprintf("Key %s does not exist", token))
			result.NoAuth("当前未登录", c)
			c.Abort()
			return
		}
		// parseToken 解析token包含的信息
		claims, err := m.Jwt.ParseToken(token)
		if err == nil {
			// 需要刷新token
			// 登出请求不需要刷新token
			if jwt.IsWithinRefreshWindow(claims) && c.FullPath() != "/api/auth/logout" {
				result.ShouldRefresh(c)
				c.Abort()
				return
			}
			c.Set("claims", claims)
			c.Next()
			return
		}
		// token 过期
		if errors.Is(err, jwt.TokenExpired) {
			result.NoAuth("授权已过期", c)
			c.Abort()
			return
		}
		// token格式无效或者错误
		if errors.Is(err, jwt.TokenNotValidYet) ||
			errors.Is(err, jwt.TokenMalformed) ||
			errors.Is(err, jwt.TokenInvalid) {
			result.NoAuth("token无效", c)
			c.Abort()
			return
		}
	}
}

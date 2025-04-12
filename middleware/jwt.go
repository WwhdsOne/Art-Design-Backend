package middleware

import (
	"Art-Design-Backend/global"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/redisx"
	"Art-Design-Backend/pkg/response"
	"Art-Design-Backend/pkg/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := utils.GetToken(c)
		id := redisx.Get(constant.LOGIN + token)
		if id == "" {
			global.Logger.Error(fmt.Sprintf("Key %s does not exist", token))
			response.NoAuth("token失效", c)
			c.Abort()
			return
		} else {
			global.Logger.Info(fmt.Sprintf("Value for key %s", token))
		}
		// parseToken 解析token包含的信息
		claims, err := global.JWT.ParseToken(token)
		if err != nil {
			if errors.Is(err, jwt.TokenExpired) {
				response.NoAuth("授权已过期", c)
				c.Abort()
				return
			}
			global.Logger.Error(fmt.Sprintf("Token解析失败：%v", err))
			response.NoAuth(err.Error(), c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

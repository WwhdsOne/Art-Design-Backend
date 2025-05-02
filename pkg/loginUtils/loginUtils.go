package loginUtils

import (
	"Art-Design-Backend/pkg/jwt"
	"context"
	"github.com/gin-gonic/gin"
)

func GetUserID(c context.Context) int64 {
	value := c.Value("claims")
	if value != nil {
		return value.(*jwt.CustomClaims).BaseClaims.ID
	}
	// 不存在操作用户则返回id为-1
	return -1
}

// GetToken 从header中获取authorization
func GetToken(c *gin.Context) string {
	token := c.GetHeader("authorization")
	return token
}

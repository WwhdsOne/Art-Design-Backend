package authutils

import (
	"Art-Design-Backend/pkg/jwt"
	"context"
	"github.com/gin-gonic/gin"
)

func GetUserID(c context.Context) int64 {
	value := c.Value("claims")
	if value != nil {
		return value.(*jwt.CustomClaims).BaseClaims.UserID
	}
	// 不存在操作用户则返回id为-1
	return -1
}

func GetUserRoleIDs(c context.Context) (roleIds []int64) {
	value := c.Value("claims")
	if value != nil {
		return value.(*jwt.CustomClaims).BaseClaims.RoleIDs
	}
	// 不存在操作用户则返回id为-1
	return
}

func GetClaims(c context.Context) *jwt.CustomClaims {
	value := c.Value("claims")
	if value != nil {
		return value.(*jwt.CustomClaims)
	}
	return nil
}

// GetToken 从header中获取authorization
func GetToken(c *gin.Context) string {
	token := c.GetHeader("authorization")
	return token
}

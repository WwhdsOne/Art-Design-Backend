package initialize

import (
	"Art-Design-Backend/api"
	"Art-Design-Backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	allGroup := r.Group("/api")
	{
		// 公共路由组（无需认证）
		openAPIGroup := allGroup.Group("/")
		api.InitOpenAuthRouter(openAPIGroup)
	}
	{
		// 私有路由组（需要 JWT 认证）
		securedAPIGroup := allGroup.Group("/")
		securedAPIGroup.Use(middleware.JWTAuth())
		api.InitSecuredAuthRouter(securedAPIGroup)
		api.InitUserRouter(securedAPIGroup)
	}

}

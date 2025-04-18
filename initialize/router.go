package initialize

import (
	"Art-Design-Backend/api"
	"Art-Design-Backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	allGroup := r.Group("/api")
	{
		publicGroup := allGroup.Group("/")
		api.InitAuthRouter(publicGroup)
	}
	{
		privateGroup := allGroup.Group("/")
		privateGroup.Use(middleware.JWTAuth())
		api.InitUserRouter(privateGroup)
	}

}

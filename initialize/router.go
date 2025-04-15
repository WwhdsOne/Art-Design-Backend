package initialize

import (
	"Art-Design-Backend/api"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	allGroup := r.Group("/api/v1")
	{
		publicGroup := allGroup.Group("/")
		api.InitAuthRouter(publicGroup)
	}
	{
		privateGroup := allGroup.Group("/")
		//privateGroup.Use(middleware.JWTAuth())
		api.InitUserRouter(privateGroup)
	}

}

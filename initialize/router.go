package initialize

import (
	"Art-Design-Backend/api"
	"Art-Design-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	{
		publicGroup := r.Group("/")
		api.InitAuthRouter(publicGroup)
	}
	{
		privateGroup := r.Group("/")
		privateGroup.Use(middleware.JWTAuth())
		api.InitUserRouter(privateGroup)
	}

}

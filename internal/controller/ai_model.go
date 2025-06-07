package controller

import (
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type AIModelController struct {
	aiModelService *service.AIModelService // 创建一个AIModelService实例
}

func NewAIModelController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.AIModelService) *AIModelController {
	aiModelCtrl := &AIModelController{
		aiModelService: service,
	}
	_ = engine.Group("/api").Group("/aimodel").Use(middleware.AuthMiddleware())

	return aiModelCtrl
}

package controller

import (
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"
	"github.com/gin-gonic/gin"
)

type AIModelController struct {
	aiModelService *service.AIModelService // 创建一个AIModelService实例
}

func NewAIModelController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.AIModelService) *AIModelController {
	aiModelCtrl := &AIModelController{
		aiModelService: service,
	}
	r := engine.Group("/api").Group("/aimodel").Use(middleware.AuthMiddleware())
	{
		r.POST("/create", aiModelCtrl.createAIModel)
		r.POST("/page", aiModelCtrl.getAIModelPage)
		r.POST("/chat-completion", aiModelCtrl.chatCompletion)
		r.GET("/getSimpleModelList", aiModelCtrl.getSimpleModelList)
	}
	return aiModelCtrl
}

func (a *AIModelController) createAIModel(c *gin.Context) {
	var req request.AIModel
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := a.aiModelService.CreateAIModel(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("模型创建成功", c)
}

func (a *AIModelController) getAIModelPage(c *gin.Context) {
	var req query.AIModel
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	res, err := a.aiModelService.GetAIModelPage(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

func (a *AIModelController) chatCompletion(c *gin.Context) {
	var req request.ChatCompletion
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := a.aiModelService.ChatCompletion(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
}

func (a *AIModelController) getSimpleModelList(c *gin.Context) {
	res, err := a.aiModelService.GetSimpleModelList(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

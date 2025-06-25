package controller

import (
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"
	"github.com/gin-gonic/gin"
)

type AIController struct {
	aiService *service.AIService // 创建一个AIService实例
}

func NewAIController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.AIService) *AIController {
	aiModelCtrl := &AIController{
		aiService: service,
	}
	r := engine.Group("/api").Group("/ai")
	{
		aiModelGroup := r.Group("/model")
		aiModelGroup.Use(middleware.AuthMiddleware())
		aiModelGroup.POST("/create", aiModelCtrl.createAIModel)
		aiModelGroup.POST("/page", aiModelCtrl.getAIModelPage)
		//r.POST("/chat-completion", aiModelCtrl.chatCompletion)
		aiModelGroup.GET("/getSimpleModelList", aiModelCtrl.getSimpleModelList)
	}
	{
		aiProviderGroup := r.Group("/provider")
		aiProviderGroup.Use(middleware.AuthMiddleware())
		aiProviderGroup.POST("/create", aiModelCtrl.createAIProvider)
		aiProviderGroup.POST("/page", aiModelCtrl.getAIProviderPage)
	}
	return aiModelCtrl
}

func (a *AIController) createAIProvider(c *gin.Context) {
	var req request.AIProvider
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := a.aiService.CreateAIProvider(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("供应商创建成功", c)
}

func (a *AIController) getAIProviderPage(c *gin.Context) {
	var q query.AIProvider
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(err)
		return
	}
	res, err := a.aiService.GetAIProviderPage(c, &q)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

func (a *AIController) createAIModel(c *gin.Context) {
	var req request.AIModel
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := a.aiService.CreateAIModel(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("模型创建成功", c)
}

func (a *AIController) getAIModelPage(c *gin.Context) {
	var req query.AIModel
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	res, err := a.aiService.GetAIModelPage(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

//func (a *AIController) chatCompletion(c *gin.Context) {
//	var req request.ChatCompletion
//	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
//		_ = c.Error(err)
//		return
//	}
//	if err := a.aiService.ChatCompletion(c, &req); err != nil {
//		_ = c.Error(err)
//		return
//	}
//}

func (a *AIController) getSimpleModelList(c *gin.Context) {
	res, err := a.aiService.GetSimpleModelList(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

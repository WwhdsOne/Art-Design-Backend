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
	aiCtrl := &AIController{
		aiService: service,
	}
	r := engine.Group("/api").Group("/ai")
	{
		aiModelGroup := r.Group("/model")
		aiModelGroup.Use(middleware.AuthMiddleware())
		aiModelGroup.POST("/create", aiCtrl.createAIModel)
		aiModelGroup.POST("/page", aiCtrl.getAIModelPage)
		aiModelGroup.POST("/chat-completion", aiCtrl.chatCompletion)
		aiModelGroup.GET("/getSimpleModelList", aiCtrl.getSimpleModelList)
		aiModelGroup.POST("/uploadIcon", aiCtrl.uploadAIModelIcon)
	}
	{
		aiProviderGroup := r.Group("/provider")
		aiProviderGroup.Use(middleware.AuthMiddleware())
		aiProviderGroup.POST("/create", aiCtrl.createAIProvider)
		aiProviderGroup.POST("/page", aiCtrl.getAIProviderPage)
		aiProviderGroup.GET("/simpleList", aiCtrl.getSimpleProviderList)
	}
	return aiCtrl
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

func (a *AIController) uploadAIModelIcon(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		result.FailWithMessage("请选择要上传的文件", c)
		return
	}

	// 打开上传的文件流
	src, err := file.Open()
	if err != nil {
		result.FailWithMessage("无法打开上传的文件", c)
		return
	}
	defer src.Close()

	// 检查文件大小是否超过 2MB
	if file.Size > 2<<20 { // 2 MB
		result.FailWithMessage("文件大小不能超过 1MB", c)
		return
	}

	modelIconURL, err := a.aiService.UploadModelIcon(c, file.Filename, src)
	if err != nil {
		result.FailWithMessage("模型图标上传失败: "+err.Error(), c)
		return
	}

	result.OkWithData(modelIconURL, c)
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

func (a *AIController) chatCompletion(c *gin.Context) {
	var req request.ChatCompletion
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := a.aiService.ChatCompletion(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
}

func (a *AIController) getSimpleModelList(c *gin.Context) {
	res, err := a.aiService.GetSimpleModelList(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

func (a *AIController) getSimpleProviderList(c *gin.Context) {
	res, err := a.aiService.GetSimpleProviderList(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

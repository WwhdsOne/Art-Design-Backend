package controller

import (
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"
	"Art-Design-Backend/pkg/utils"

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
		aiModelGroup.GET("/simpleList", aiCtrl.getSimpleModelList)
		aiModelGroup.POST("/uploadIcon", aiCtrl.uploadAIModelIcon)
		aiModelGroup.POST("/uploadChatFile", aiCtrl.uploadChatImage)
	}
	{
		aiProviderGroup := r.Group("/provider")
		aiProviderGroup.Use(middleware.AuthMiddleware())
		aiProviderGroup.POST("/create", aiCtrl.createAIProvider)
		aiProviderGroup.POST("/page", aiCtrl.getAIProviderPage)
		aiProviderGroup.GET("/simpleList", aiCtrl.getSimpleProviderList)
	}
	{
		//agentGroup := r.Group("/agent")
		//agentGroup.Use(middleware.AuthMiddleware())
		//agentGroup.POST("/uploadAgentDocument/:id", aiCtrl.UploadAgentDocument)
		//agentGroup.POST("/create", aiCtrl.CreateAgent)
		//agentGroup.POST("/page", aiCtrl.GetAIAgentPage)
		//agentGroup.GET("/simpleList", aiCtrl.getSimpleAgentList)
		//agentGroup.POST("/chat-completion", aiCtrl.agentChatCompletion)
	}
	{
		conversationGroup := r.Group("/conversation")
		conversationGroup.Use(middleware.AuthMiddleware())
		conversationGroup.GET("/history", aiCtrl.GetHistoryConversation)
		conversationGroup.GET("/:id/messages", aiCtrl.GetMessageByConversationID)
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
	res, err := a.aiService.GetSimpleChatModelList(c)
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

func (a *AIController) GetHistoryConversation(c *gin.Context) {
	res, err := a.aiService.GetHistoryConversation(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

func (a *AIController) GetMessageByConversationID(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	res, err := a.aiService.GetMessageByConversationID(c, id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

func (a *AIController) uploadChatImage(c *gin.Context) {
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

	// 检查文件大小是否超过 10MB
	if file.Size > 10<<20 { // 10 MB
		result.FailWithMessage("文件大小不能超过 10MB", c)
		return
	}

	chatMessageURL, err := a.aiService.UploadChatMessageImage(c, file.Filename, src)
	if err != nil {
		result.FailWithMessage("模型图标上传失败: "+err.Error(), c)
		return
	}

	result.OkWithData(chatMessageURL, c)
}

//func (a *AIController) UploadAgentDocument(c *gin.Context) {
//	file, header, err := c.Request.FormFile("file")
//	if err != nil {
//		result.FailWithMessage("请选择要上传的文件", c)
//		return
//	}
//	defer file.Close()
//
//	if header.Size > 100<<20 {
//		result.FailWithMessage("文件不能超过 100MB", c)
//		return
//	}
//
//	filename := header.Filename
//	agentID, err := utils.ParseID(c)
//
//	err = a.aiService.UploadAndVectorizeDocument(c, file, filename, agentID)
//	if err != nil {
//		result.FailWithMessage("文件上传失败: "+err.Error(), c)
//		return
//	}
//
//	result.OkWithMessage("上传成功", c)
//}

//func (a *AIController) CreateAgent(c *gin.Context) {
//	var req request.AIAgent
//	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
//		_ = c.Error(err)
//		return
//	}
//	if err := a.aiService.CreateAgent(c, &req); err != nil {
//		_ = c.Error(err)
//		return
//	}
//	result.OkWithMessage("创建成功", c)
//}
//
//func (a *AIController) getSimpleAgentList(c *gin.Context) {
//	res, err := a.aiService.GetSimpleAgentList(c)
//	if err != nil {
//		_ = c.Error(err)
//		return
//	}
//	result.OkWithData(res, c)
//}
//
//func (a *AIController) GetAIAgentPage(c *gin.Context) {
//	var req query.AIAgent
//	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
//		_ = c.Error(err)
//		return
//	}
//	res, err := a.aiService.GetAIAgentPage(c, &req)
//	if err != nil {
//		_ = c.Error(err)
//		return
//	}
//	result.OkWithData(res, c)
//}

//func (a *AIController) agentChatCompletion(c *gin.Context) {
//	var req request.ChatCompletion
//	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
//		_ = c.Error(err)
//		return
//	}
//	if err := a.aiService.AgentChatCompletion(c, &req); err != nil {
//		_ = c.Error(err)
//		return
//	}
//}

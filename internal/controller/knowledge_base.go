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

type KnowledgeBaseController struct {
	knowledgeBaseService *service.KnowledgeBaseService
}

func NewKnowledgeBaseController(engine *gin.Engine, mws *middleware.Middlewares, svc *service.KnowledgeBaseService) *KnowledgeBaseController {
	knowledgeBaseCtrl := &KnowledgeBaseController{
		knowledgeBaseService: svc,
	}
	r := engine.Group("/api").Group("/knowledgeBase")
	{
		r.Use(mws.AuthMiddleware())
		r.POST("/file/page", knowledgeBaseCtrl.GetKnowledgeBaseFileList)
		r.POST("/file/upload", knowledgeBaseCtrl.UploadFile)
		//r.GET("/list", knowledgeBaseCtrl.GetKnowledgeBaseFileList)
	}
	{
		r.POST("/page", knowledgeBaseCtrl.GetKnowledgeBasePage)
		r.POST("/create", knowledgeBaseCtrl.CreateKnowledgeBase)
		r.POST("/delete/:id", knowledgeBaseCtrl.DeleteKnowledgeBase)
		r.POST("/update", knowledgeBaseCtrl.UpdateKnowledgeBase)
		r.GET("/:id/files", knowledgeBaseCtrl.GetKnowledgeBaseFilesByID)
		r.GET("/simpleList", knowledgeBaseCtrl.GetSimpleKnowledgeBaseList)
	}
	return knowledgeBaseCtrl
}

func (k *KnowledgeBaseController) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		result.FailWithMessage("请选择要上传的文件", c)
		return
	}
	defer file.Close()

	fileSize := header.Size

	if fileSize > 100<<20 {
		result.FailWithMessage("文件不能超过 100MB", c)
		return
	}

	filename := header.Filename

	err = k.knowledgeBaseService.UploadAndVectorizeDocument(c, file, filename, fileSize)
	if err != nil {
		result.FailWithMessage("文件上传失败: "+err.Error(), c)
		return
	}

	result.OkWithMessage("上传成功", c)
}

func (k *KnowledgeBaseController) GetKnowledgeBaseFileList(c *gin.Context) {
	var search query.KnowledgeBaseFile
	// 如果绑定过程中出现错误，返回错误响应并结束函数执行
	if err := c.ShouldBindBodyWithJSON(&search); err != nil {
		_ = c.Error(err)
		return
	}
	res, err := k.knowledgeBaseService.GetKnowledgeBaseFileList(c, &search)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

func (k *KnowledgeBaseController) GetKnowledgeBasePage(c *gin.Context) {
	var search query.KnowledgeBase
	// 如果绑定过程中出现错误，返回错误响应并结束函数执行
	if err := c.ShouldBindBodyWithJSON(&search); err != nil {
		_ = c.Error(err)
		return
	}
	res, err := k.knowledgeBaseService.GetKnowledgeBasePage(c, &search)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

func (k *KnowledgeBaseController) CreateKnowledgeBase(c *gin.Context) {
	var req request.KnowledgeBase
	// 如果绑定过程中出现错误，返回错误响应并结束函数执行
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	err := k.knowledgeBaseService.CreateKnowledgeBase(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("知识库创建成功", c)
}

func (k *KnowledgeBaseController) DeleteKnowledgeBase(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = k.knowledgeBaseService.DeleteKnowledgeBase(c, id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("知识库删除成功", c)
}

func (k *KnowledgeBaseController) UpdateKnowledgeBase(c *gin.Context) {
	var req request.KnowledgeBase
	// 如果绑定过程中出现错误，返回错误响应并结束函数执行
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	err := k.knowledgeBaseService.UpdateKnowledgeBase(c, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("知识库更新成功", c)
}

func (k *KnowledgeBaseController) GetKnowledgeBaseFilesByID(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	res, err := k.knowledgeBaseService.GetKnowledgeBaseFilesByID(c, id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

func (k *KnowledgeBaseController) GetSimpleKnowledgeBaseList(c *gin.Context) {
	res, err := k.knowledgeBaseService.GetSimpleKnowledgeBaseList(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

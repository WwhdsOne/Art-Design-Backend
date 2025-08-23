package controller

import (
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"

	"github.com/gin-gonic/gin"
)

type KnowledgeBaseController struct {
	knowledgeBaseService *service.KnowledgeBaseService
}

func NewKnowledgeBaseController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.KnowledgeBaseService) *KnowledgeBaseController {
	knowledgeBaseCtrl := &KnowledgeBaseController{
		knowledgeBaseService: service,
	}
	r := engine.Group("/api").Group("/knowledgeBase")
	{
		r.Use(middleware.AuthMiddleware())
		r.POST("/uploadFile", knowledgeBaseCtrl.UploadFile)
		//r.GET("/list", knowledgeBaseCtrl.GetKnowledgeBaseFileList)
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

package controller

import (
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DigitPredictController struct {
	DigitPredictService *service.DigitPredictService
}

func NewDigitPredictController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.DigitPredictService) *DigitPredictController {
	digitCtrl := &DigitPredictController{
		DigitPredictService: service,
	}
	r := engine.Group("/api").Group("/digitPredict")
	r.Use(middleware.AuthMiddleware())
	{
		// 私有路由组（需要 JWT 认证）
		r.GET("/list", digitCtrl.getDigitPredictList)
		r.POST("/predict", digitCtrl.predict)
		r.POST("/upload", digitCtrl.uploadDigitImage)
	}
	return digitCtrl
}

func (d *DigitPredictController) getDigitPredictList(c *gin.Context) {
	userID := authutils.GetUserID(c)
	list, err := d.DigitPredictService.GetDigitPredictList(c, userID)
	if err != nil {
		result.FailWithMessage("获取数字识别列表失败", c)
		return
	}
	result.OkWithData(list, c)
}

func (d *DigitPredictController) predict(c *gin.Context) {
	var req request.DigitPredict
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		c.Set(gin.BindKey, req)
		return
	}
	_ = d.DigitPredictService.SubmitMission(c, &req)
	result.OkWithMessage("提交成功", c)
}

func (d *DigitPredictController) uploadDigitImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("上传文件为空", zap.Error(err))
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

	// 检查文件大小是否超过 8MB
	if file.Size > 2<<20 { // 2 MB
		result.FailWithMessage("文件大小不能超过 2MB", c)
		return
	}

	// 调用 service 层处理上传逻辑
	digitImageUrl, err := d.DigitPredictService.UploadDigitImage(c, file.Filename, src)
	if err != nil {
		result.FailWithMessage("数字上传失败: "+err.Error(), c)
		return
	}

	result.OkWithData(digitImageUrl, c)
}

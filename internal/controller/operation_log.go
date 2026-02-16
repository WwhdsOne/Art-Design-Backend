package controller

import (
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"

	"github.com/gin-gonic/gin"
)

type OperationLogController struct {
	operationLogService *service.OperationLogService
}

func NewOperationLogController(engine *gin.Engine, mws *middleware.Middlewares, svc *service.OperationLogService) *OperationLogController {
	operationLogCtrl := &OperationLogController{
		operationLogService: svc,
	}
	r := engine.Group("/api").Group("/operationLog")
	r.Use(mws.AuthMiddleware())
	r.POST("/page", operationLogCtrl.getOperationPage)
	return operationLogCtrl
}

func (o *OperationLogController) getOperationPage(c *gin.Context) {
	var operationQuery query.OperationLog
	if err := c.ShouldBindBodyWithJSON(&operationQuery); err != nil {
		_ = c.Error(err)
		return
	}
	operationLogPageResp, err := o.operationLogService.GetOperationLogPage(c, &operationQuery)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(operationLogPageResp, c)
}

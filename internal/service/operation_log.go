package service

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type OperationLogService struct {
	OperationLogRepo *repository.OperationLogRepo
}

func (o *OperationLogService) GetOperationLogPage(
	c *gin.Context,
	logQuery *query.OperationLog,
) (resp *common.PaginationResp[response.OperationLog], err error) {

	// entity 层数据
	logEntities, total, err := o.OperationLogRepo.GetOperationLogPage(c, logQuery)
	if err != nil {
		return
	}

	// response 层数据
	logRespList := make([]response.OperationLog, 0, len(logEntities))
	for _, logEntity := range logEntities {
		var logResp response.OperationLog
		_ = copier.Copy(&logResp, logEntity)
		logRespList = append(logRespList, logResp)
	}

	resp = common.BuildPageResp[response.OperationLog](
		logRespList,
		total,
		logQuery.PaginationReq,
	)

	return
}

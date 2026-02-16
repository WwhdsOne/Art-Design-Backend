package query

import (
	"Art-Design-Backend/internal/model/common"
	"time"
)

type OperationLog struct {
	OperatorID int64  `json:"operator_id"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Status     int16  `json:"status"`
	IP         string `json:"ip"`
	Browser    string `json:"browser"`
	OS         string `json:"os"`

	// ===== 时间区间查询 =====
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`

	common.PaginationReq
}

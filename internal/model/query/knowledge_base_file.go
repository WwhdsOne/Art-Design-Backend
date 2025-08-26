package query

import (
	"Art-Design-Backend/internal/model/common"
	"time"
)

type KnowledgeBaseFile struct {
	common.PaginationReq
	OriginalFileName *string     `json:"original_filename"`    // 模糊匹配文件名
	FileType         *string     `json:"file_type"`            // 支持多文件类型过滤
	CreateUser       *string     `json:"create_user"`          // 创建人
	TimeRange        []time.Time `json:"time_range,omitempty"` // 时间范围
}

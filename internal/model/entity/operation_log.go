package entity

import (
	"Art-Design-Backend/pkg/constant/tablename"
	"time"
)

type OperationLog struct {
	ID         int64     `gorm:"type:bigint;primaryKey;comment:雪花ID"` // 雪花ID
	OperatorID int64     `gorm:"type:bigint;not null;index:idx_operator_id;comment:操作人ID,非鉴权接口则为-1"`
	Method     string    `gorm:"type:varchar(10);not null;comment:HTTP请求方法(GET/POST等)"`
	Path       string    `gorm:"type:varchar(255);not null;comment:请求路径"`
	IP         string    `gorm:"type:inet;comment:客户端IP地址"`
	Params     string    `gorm:"type:text;comment:请求参数(JSON格式存储)"`
	Status     int16     `gorm:"type:smallint;not null;comment:HTTP响应状态码"`
	CreatedAt  time.Time `gorm:"type:timestamp with time zone;not null;index:idx_created_at;comment:操作时间"`
}

func (o *OperationLog) TableName() string {
	return tablename.OperationTableName
}

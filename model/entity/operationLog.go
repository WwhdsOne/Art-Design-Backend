package entity

import (
	"Art-Design-Backend/model/base"
	"github.com/dromara/carbon/v2"
)

type OperationLog struct {
	base.ID
	UserID    int64           `gorm:"column:user_id;type:bigint;not null;index;comment:'操作人ID'"`
	Method    string          `gorm:"column:method;type:varchar(10);not null;comment:'HTTP请求方法(GET/POST等)'"`
	Path      string          `gorm:"column:path;type:varchar(255);not null;comment:'请求路径'"`
	IP        string          `gorm:"column:ip;type:varchar(50);not null;comment:'客户端IP地址'"`
	Params    string          `gorm:"column:params;type:text;comment:'请求参数(JSON格式存储)'"`
	Status    int16           `gorm:"column:status;type:smallint;not null;comment:'HTTP响应状态码'"`
	CreatedAt carbon.DateTime `gorm:"column:created_at;type:timestamp;not null;index;comment:'操作时间'"`
}

func (o *OperationLog) TableName() string {
	return "operation_log"
}

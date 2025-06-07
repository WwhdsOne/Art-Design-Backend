package entity

import (
	"Art-Design-Backend/pkg/constant/tablename"
	"time"
)

type OperationLog struct {
	ID             int64     `gorm:"type:bigint;primaryKey;comment:雪花ID"`
	OperatorID     int64     `gorm:"type:bigint;not null;index:idx_operator_id;comment:操作人ID,非鉴权接口则为-1"`
	Method         string    `gorm:"type:varchar(10);not null;comment:HTTP请求方法(GET/POST等)"`
	Path           string    `gorm:"type:varchar(255);not null;comment:请求路径"`
	Query          string    `gorm:"type:text;comment:URL查询参数"`
	Params         string    `gorm:"type:text;comment:请求参数(JSON格式存储)"`
	Status         int16     `gorm:"type:smallint;not null;comment:HTTP响应状态码"`
	Latency        int64     `gorm:"comment:请求耗时(ms)"`
	IP             string    `gorm:"type:varchar(64);comment:客户端IP地址"` // `inet` 可兼容性差，建议用 varchar
	UserAgent      string    `gorm:"type:varchar(255);comment:浏览器或客户端信息"`
	Browser        string    `gorm:"type:varchar(50);comment:浏览器名称"`
	BrowserVersion string    `gorm:"type:varchar(50);comment:浏览器版本"`
	OS             string    `gorm:"type:varchar(50);comment:操作系统名称"`
	Platform       string    `gorm:"type:varchar(50);comment:平台名称"`
	ErrorMsg       string    `gorm:"type:text;comment:错误信息"`
	CreatedAt      time.Time `gorm:"type:timestamp;not null;index:idx_created_at;comment:操作时间"`
}

func (o *OperationLog) TableName() string {
	return tablename.OperationTableName
}

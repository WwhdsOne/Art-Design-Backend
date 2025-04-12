package entity

import "github.com/dromara/carbon/v2"

type OperationLog struct {
	ID        int64           `gorm:"primarykey;column:id;comment:雪花ID" json:"id"`                        // 雪花ID
	UserID    int64           `gorm:"column:user_id;comment:操作人ID" json:"user_id"`                        // 操作人 ID
	Method    string          `gorm:"type:varchar(30);column:method;comment:HTTP方法" json:"method"`        // HTTP 方法
	Path      string          `gorm:"type:varchar(255);column:path;comment:请求路径" json:"path"`             // 请求路径
	IP        string          `gorm:"type:varchar(255);column:ip;comment:客户端IP" json:"ip"`                // 客户端 IP
	Params    string          `gorm:"type:varchar(512);column:params;comment:请求参数(JSON格式)" json:"params"` // 请求参数（JSON 格式）
	Status    int16           `gorm:"type:smallint;column:status;comment:HTTP状态码" json:"status"`          // HTTP 状态码
	CreatedAt carbon.DateTime `gorm:"comment:创建时间" json:"created_at"`                                     // 创建时间
}

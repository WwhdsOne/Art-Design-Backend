package resp

import "github.com/dromara/carbon/v2"

type User struct {
	ID             int64           `json:"id,string"`
	Username       string          `json:"username"`
	Nickname       string          `json:"nickname"`
	CreatedAt      carbon.DateTime `json:"createdAt"` // 创建记录时自动设置为当前时间
	CreateBy       int64           `json:"createBy"`  // 创建人字段，记录创建操作者的标识
	CreateUserName string          `gorm:"-" json:"createUserName"`
}

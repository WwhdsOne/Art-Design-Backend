package entity

import (
	"Art-Design-Backend/model/base"
)

// Permission 定义权限模型
type Permission struct {
	base.BaseModel
	Name     string  `gorm:"varchar(20);unique;not null;comment:'页面名称'"` // 页面名称
	ParentId int64   `gorm:"int8;not null;comment:'父页面ID'"`              // 父页面ID
	Path     string  `gorm:"varchar(256);comment:'路径'"`                  // 路径
	ChildIds []int64 `gorm:"-"`                                          // 子页面ID列表，不直接存储在数据库中
}

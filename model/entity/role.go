package entity

import (
	"Art-Design-Backend/model/base"
)

// Role 定义角色模型
type Role struct {
	base.BaseModel
	Name        string        `gorm:"varchar(10);unique;not null"`
	Status      int8          `gorm:"type:tinyint;not null;default:1;comment:'状态，1表示正常，0表示禁用'"`
	Description string        `gorm:"type:varchar(256);comment:'描述'"`
	Permissions []*Permission `gorm:"many2many:role_permissions;"`
}

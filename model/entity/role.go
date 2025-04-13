package entity

import "Art-Design-Backend/model/base"

// Role 定义角色模型
type Role struct {
	base.BaseModel
	Name        string       `gorm:"column:name;type:varchar(10);not null;unique;comment:'角色名称(唯一)'"`
	Description string       `gorm:"column:description;type:varchar(256);comment:'角色描述'"`
	Status      int8         `gorm:"column:status;type:tinyint;not null;default:1;comment:'状态:1-正常,0-禁用'"`
	Permissions []Permission `gorm:"many2many:role_permissions;comment:'关联权限'"`
}

func (r *Role) TableName() string {
	return "role"
}

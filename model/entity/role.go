package entity

import "Art-Design-Backend/model/base"

// Role 定义角色模型
type Role struct {
	base.BaseModel
	Name        string `gorm:"type:varchar(30);not null;unique;comment:角色名称"`
	Code        string `gorm:"type:varchar(30);not null;unique;comment:角色编码"`
	Description string `gorm:"type:varchar(256);comment:角色描述"`
	Status      int8   `gorm:"type:tinyint;not null;default:1;comment:状态:0-禁用,1-正常"`
	Menus       []Menu `gorm:"many2many:role_menus;comment:关联权限"`
}

func (r *Role) TableName() string {
	return "role"
}

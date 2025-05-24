package entity

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/pkg/constant/tablename"
)

// Role 定义角色模型
type Role struct {
	base.BaseModel
	Name        string `gorm:"type:varchar(30);not null;unique;comment:角色名称"`
	Code        string `gorm:"type:varchar(30);not null;unique;comment:角色编码"`
	Description string `gorm:"type:varchar(256);comment:角色描述"`
	Status      int8   `gorm:"type:smallint;not null;default:1;comment:状态:1-正常,2-禁用"`
	Menus       []Menu `gorm:"many2many:role_menus;comment:关联权限" json:"menus,omitzero"`
}

func (r *Role) TableName() string {
	return tablename.RoleTableName
}

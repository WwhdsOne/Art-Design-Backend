package entity

import "Art-Design-Backend/pkg/constant"

// RoleMenus 角色-菜单关联表
type RoleMenus struct {
	RoleID int64
	MenuID int64
}

func (r *RoleMenus) TableName() string {
	return constant.RoleMenusTableName
}

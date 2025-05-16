package entity

import "Art-Design-Backend/pkg/constant/tablename"

// UserRoles 定义用户和角色的多对多关联表
type UserRoles struct {
	UserID int64
	RoleID int64
}

func (u *UserRoles) TableName() string {
	return tablename.UserRolesTableName
}

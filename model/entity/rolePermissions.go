package entity

// RolePermissions 角色-权限关联表
type RolePermissions struct {
	RoleID       int64 `gorm:"column:role_id;type:bigint;not null;index;comment:'角色ID'"`
	PermissionID int64 `gorm:"column:permission_id;type:bigint;not null;comment:'权限ID'"`
}

func (r *RolePermissions) TableName() string {
	return "role_permissions"
}

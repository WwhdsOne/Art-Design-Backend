package entity

// UserRoles 定义用户和角色的多对多关联表
type UserRoles struct {
	UserID int64
	RoleID int64
}

func (u *UserRoles) TableName() string {
	return "user_roles"
}

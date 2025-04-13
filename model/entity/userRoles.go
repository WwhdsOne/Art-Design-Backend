package entity

// UserRoles 定义用户和角色的多对多关联表
type UserRoles struct {
	UserID int64 `gorm:"column:user_id;type:bigint;not null;comment:用户ID"` // 用户ID，非空且建立索引
	RoleID int64 `gorm:"column:role_id;type:bigint;not null;comment:角色ID"` // 角色ID，非空且建立索引
}

func (u *UserRoles) TableName() string {
	return "user_roles"
}

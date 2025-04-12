package entity

// UserRoles 定义用户和角色的多对多关联表
type UserRoles struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

package entity

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/pkg/constant"
)

// User 结构体定义了用户的基本信息
type User struct {
	base.BaseModel
	Username     string   `gorm:"type:varchar(32);uniqueIndex;not null;comment:用户名"`
	RealName     string   `gorm:"type:varchar(50);comment:真实姓名"`
	Nickname     string   `gorm:"type:varchar(24);not null;comment:昵称"`
	Password     string   `gorm:"type:varchar(255);not null;comment:密码（加密存储）"`
	Gender       int8     `gorm:"type:tinyint;not null;default:1;comment:性别:1-男,2-女"`
	Email        *string  `gorm:"type:varchar(256);uniqueIndex;comment:邮箱" copier:"-"`
	Phone        *string  `gorm:"type:varchar(256);uniqueIndex;comment:手机号" copier:"-"`
	Address      string   `gorm:"type:varchar(256);comment:地址"`
	Avatar       string   `gorm:"type:varchar(255);comment:头像URL"`
	Introduction string   `gorm:"type:varchar(256);comment:个人介绍"`
	Occupation   string   `gorm:"type:varchar(50);comment:职业"`
	Tags         []string `gorm:"type:json;serializer:json;comment:个人标签"`
	Status       int8     `gorm:"type:tinyint;not null;default:1;comment:状态:0-禁用,1-正常"`
	Roles        []Role   `gorm:"many2many:user_roles;comment:关联角色"`
}

func (u *User) TableName() string {
	return constant.UserTableName
}

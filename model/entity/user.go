package entity

import (
	"Art-Design-Backend/model/base"
)

// User 结构体定义了用户的基本信息
type User struct {
	base.BaseModel
	Username     string   `gorm:"column:username;type:varchar(32);unique_index;not null;comment:'用户名'"`
	RealName     string   `gorm:"column:real_name;type:varchar(20);comment:'真实姓名'"`
	Nickname     string   `gorm:"column:nickname;type:varchar(12);not null;comment:'昵称'"`
	Password     string   `gorm:"column:password;type:varchar(32);not null;comment:'密码'"`
	Gender       int8     `gorm:"column:gender;type:tinyint;not null;default:1;comment:'性别:1表示男性，2表示女性'"`
	Email        string   `gorm:"column:email;type:varchar(100);unique_index;comment:'邮箱地址，唯一索引'"`
	Phone        string   `gorm:"column:phone;type:varchar(30);unique_index;comment:'手机号码，唯一索引'"`
	Address      string   `gorm:"column:address;type:varchar(256);comment:'地址'"`
	Introduction string   `gorm:"column:introduction;type:varchar(256);comment:'个人介绍'"`
	Occupation   string   `gorm:"column:occupation;type:varchar(50);comment:'职业'"`
	Tags         []string `gorm:"column:tags;type:json;comment:'个人标签'"`
	Status       int8     `gorm:"column:status;type:tinyint;not null;default:1;comment:'状态，1表示正常，0表示禁用'"`
	Roles        []Role   `gorm:"many2many:user_roles;comment:'关联角色'"`
}

func (u *User) TableName() string {
	return "user"
}

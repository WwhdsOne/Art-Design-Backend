package entity

import (
	"Art-Design-Backend/model/base"
	"Art-Design-Backend/model/request"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 结构体定义了用户的基本信息
type User struct {
	base.BaseModel
	Username     string   `gorm:"column:username;type:varchar(32);uniqueIndex;not null;comment:'用户名'"`
	RealName     string   `gorm:"column:real_name;type:varchar(50);comment:'真实姓名'"`
	Nickname     string   `gorm:"column:nickname;type:varchar(24);not null;comment:'昵称'"`
	Password     string   `gorm:"column:password;type:varchar(255);not null;comment:'密码（加密存储）'"`
	Gender       int8     `gorm:"column:gender;type:tinyint;not null;default:1;comment:'性别:1-男,2-女'"`
	Email        string   `gorm:"column:email;type:varchar(256);uniqueIndex;comment:'邮箱'"`
	Phone        string   `gorm:"column:phone;type:varchar(256);uniqueIndex;comment:'手机号'"`
	Address      string   `gorm:"column:address;type:varchar(256);comment:'地址'"`
	Avatar       string   `gorm:"column:avatar;type:varchar(255);comment:'头像URL'"`
	Introduction string   `gorm:"column:introduction;type:varchar(256);comment:'个人介绍'"`
	Occupation   string   `gorm:"column:occupation;type:varchar(50);comment:'职业'"`
	Tags         []string `gorm:"column:tags;type:json;serializer:json;comment:'个人标签'"`
	Status       int8     `gorm:"column:status;type:tinyint;not null;default:1;comment:'状态:0-禁用,1-正常'"`
	Roles        []Role   `gorm:"many2many:user_roles;comment:'关联角色'"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	var result struct {
		UsernameExists bool
		EmailExists    bool
		PhoneExists    bool
	}

	// 检查当前记录是否有ID，如果有，则在查询中排除它
	excludeID := ""
	if u.ID.Val != 0 {
		excludeID = fmt.Sprintf("AND id != %d", u.ID)
	}

	// 单次查询检查所有字段，排除当前ID
	db.Raw("SELECT"+
		"EXISTS(SELECT 1 FROM user WHERE username = ?"+excludeID+") AS username_exists",
		"EXISTS(SELECT 1 FROM user WHERE email = ?"+excludeID+") AS username_exists",
		"EXISTS(SELECT 1 FROM user WHERE phone = ?"+excludeID+") AS username_exists",
		u.Username, u.Email, u.Phone).Scan(&result)

	switch {
	case result.UsernameExists:
		err = fmt.Errorf("用户名重复")
	case result.EmailExists:
		err = fmt.Errorf("邮箱重复")
	case result.PhoneExists:
		err = fmt.Errorf("手机号重复")
	}
	return
}

// BeforeCopy 是 copier 的钩子函数
func (u *User) BeforeCopy(src interface{}) error {
	userReq, ok := src.(request.User)
	if !ok {
		return fmt.Errorf("源类型错误")
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	u.Password = string(password)
	email, _ := bcrypt.GenerateFromPassword([]byte(userReq.Email), bcrypt.DefaultCost)
	u.Email = string(email)
	phone, _ := bcrypt.GenerateFromPassword([]byte(userReq.Phone), bcrypt.DefaultCost)
	u.Phone = string(phone)
	return nil
}

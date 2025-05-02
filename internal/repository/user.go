package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/errorTypes"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB // 用户表数据库连接
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db.Table(constant.UserTableName),
	}
}

func (u *UserRepository) CheckUserDuplicate(user *entity.User) (err error) {

	var result struct {
		UsernameExists bool
		EmailExists    bool
		PhoneExists    bool
	}

	// 检查当前记录是否有ID，如果有，则在查询中排除它
	excludeID := ""
	if user.ID != 0 {
		excludeID = fmt.Sprintf("AND id != %d", user.ID)
	}

	// 单次查询检查所有字段，排除当前ID
	u.db.Raw("SELECT "+
		"EXISTS(SELECT 1 FROM user WHERE username = ? "+excludeID+") AS username_exists,"+
		"EXISTS(SELECT 1 FROM user WHERE email = ? "+excludeID+") AS email_exists,"+
		"EXISTS(SELECT 1 FROM user WHERE phone = ? "+excludeID+") AS phone_exists",
		user.Username, user.Email, user.Phone).Scan(&result)

	switch {
	case result.UsernameExists:
		err = errorTypes.NewGormError("用户名重复")
	case result.EmailExists:
		err = errorTypes.NewGormError("邮箱重复")
	case result.PhoneExists:
		err = errorTypes.NewGormError("手机号重复")
	}
	return
}

func (u *UserRepository) GetUserByUsername(c context.Context, username string) (user *entity.User, err error) {
	err = u.db.WithContext(c).
		Select("id", "password").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return
	}
	return
}

func (u *UserRepository) GetUserById(c context.Context, id int64) (user *entity.User, err error) {
	err = u.db.WithContext(c).
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return
	}
	return
}

func (u *UserRepository) CreateUser(c context.Context, user *entity.User) error {
	return u.db.WithContext(c).Create(user).Error
}

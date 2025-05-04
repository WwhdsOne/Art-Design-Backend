package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/errorTypes"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
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

	// 构建动态查询条件
	var query strings.Builder
	args := make([]interface{}, 0)
	conditions := make([]string, 0)

	// 只检查非空字段
	if user.Username != "" {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM user WHERE username = ? "+excludeID+") AS username_exists")
		args = append(args, user.Username)
	}

	if user.Email != nil {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM user WHERE email = ? "+excludeID+") AS email_exists")
		args = append(args, user.Email)
	}

	if user.Phone != nil {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM user WHERE phone = ? "+excludeID+") AS phone_exists")
		args = append(args, user.Phone)
	}

	// 如果没有需要检查的字段，直接返回
	if len(conditions) == 0 {
		return nil
	}

	// 构建完整查询
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(conditions, ","))

	// 执行查询
	u.db.Raw(query.String(), args...).Scan(&result)

	// 检查结果
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
		zap.L().Error("根据用户名查询用户失败", zap.Error(err))
		err = errorTypes.NewGormError("用户不存在")
		return
	}
	return
}

func (u *UserRepository) GetUserById(c context.Context, id int64) (user *entity.User, err error) {
	if err = u.db.WithContext(c).Where("id = ?", id).First(&user).Error; err != nil {
		zap.L().Error("根据用户id查询用户失败", zap.Error(err))
		err = errorTypes.NewGormError("用户不存在")
		return
	}
	return
}

func (u *UserRepository) CreateUser(c context.Context, user *entity.User) (err error) {
	if err = u.db.WithContext(c).Create(user).Error; err != nil {
		zap.L().Error("新增用户失败", zap.Error(err))
		return errorTypes.NewGormError("新增用户失败")
	}
	return err
}

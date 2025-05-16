package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/errors"
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
		db: db,
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
	var queryConditions strings.Builder
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
	queryConditions.WriteString("SELECT ")
	queryConditions.WriteString(strings.Join(conditions, ","))

	// 执行查询
	u.db.Raw(queryConditions.String(), args...).Scan(&result)

	// 检查结果
	switch {
	case result.UsernameExists:
		err = errors.NewDBError("用户名重复")
	case result.EmailExists:
		err = errors.NewDBError("邮箱重复")
	case result.PhoneExists:
		err = errors.NewDBError("手机号重复")
	}
	return
}

func (u *UserRepository) GetLoginUserByUsername(c context.Context, username string) (user *entity.User, err error) {
	if err = DB(c, u.db).
		Select("id", "password", "status").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		zap.L().Error("根据用户名查询用户失败")
		err = errors.NewDBError("用户不存在")
		return
	}
	return
}

func (u *UserRepository) GetUserById(c context.Context, id int64) (user *entity.User, err error) {
	if err = DB(c, u.db).
		Where("id = ?", id).
		First(&user).Error; err != nil {
		zap.L().Error("根据用户id查询用户失败")
		err = errors.NewDBError("用户不存在")
		return
	}
	return
}

func (u *UserRepository) CreateUser(c context.Context, user *entity.User) (err error) {
	if err = DB(c, u.db).Create(user).Error; err != nil {
		zap.L().Error("新增用户失败")
		return errors.NewDBError("新增用户失败")
	}
	return err
}

func (u *UserRepository) UpdateUser(c context.Context, user *entity.User) (err error) {
	if err = DB(c, u.db).Updates(user).Error; err != nil {
		zap.L().Error("更新用户失败")
		return errors.NewDBError("更新用户失败")
	}
	return err
}

func (u *UserRepository) GetUserPage(c context.Context, user *query.User) (userPage []*entity.User, total int64, err error) {
	db := DB(c, u.db)

	// 构建通用查询条件
	queryConditions := db.Model(&entity.User{})

	if user.RealName != "" {
		queryConditions = queryConditions.Where("real_name LIKE ?", "%"+user.RealName+"%")
	}
	if user.Username != "" {
		queryConditions = queryConditions.Where("username LIKE ?", "%"+user.Username+"%")
	}
	if user.Email != "" {
		queryConditions = queryConditions.Where("email LIKE ?", "%"+user.Email+"%")
	}
	if user.Phone != "" {
		queryConditions = queryConditions.Where("phone LIKE ?", "%"+user.Phone+"%")
	}

	if user.Gender != 0 {
		queryConditions = queryConditions.Where("gender = ?", user.Gender)
	}

	if user.Status != 0 {
		queryConditions = queryConditions.Where("status = ?", user.Status)
	}

	// 查询总数
	if err = queryConditions.Count(&total).Error; err != nil {
		zap.L().Error("获取用户分页失败")
		err = errors.NewDBError("获取用户分页失败")
		return
	}

	// 查询分页数据（可根据需要添加 Limit / Offset 支持）
	if err = queryConditions.Scopes(user.Paginate()).Find(&userPage).Error; err != nil {
		zap.L().Error("获取用户分页数据失败")
		err = errors.NewDBError("获取用户分页数据失败")
		return
	}

	return
}

package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/repository/dupcheck"
	"Art-Design-Backend/pkg/errors"
	"context"

	"gorm.io/gorm"
)

type UserDB struct {
	db *gorm.DB // 用户表数据库连接
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{
		db: db,
	}
}

func (u *UserDB) CheckUserDuplicate(c context.Context, user *entity.User) error {
	email, emailIsNil := "", user.Email == nil
	if !emailIsNil {
		email = *user.Email
	}
	phone, phoneIsNil := "", user.Phone == nil
	if !phoneIsNil {
		phone = *user.Phone
	}
	return dupcheck.Check(DB(c, u.db), "\"user\"", user.ID, []dupcheck.Field{
		{Column: "username", ErrMsg: "用户名重复", Value: user.Username},
		{Column: "email", ErrMsg: "邮箱重复", Value: email, IsNil: emailIsNil},
		{Column: "phone", ErrMsg: "手机号重复", Value: phone, IsNil: phoneIsNil},
	})
}

func (u *UserDB) GetLoginUserByUsername(c context.Context, username string) (user *entity.User, err error) {
	if err = DB(c, u.db).
		Select("id", "password", "status").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		err = errors.WrapDBError(err, "用户不存在")
		return
	}
	return
}

func (u *UserDB) GetUserByID(c context.Context, id int64) (user *entity.User, err error) {
	if err = DB(c, u.db).
		Where("id = ?", id).
		First(&user).Error; err != nil {
		err = errors.WrapDBError(err, "用户不存在")
		return
	}
	return
}

func (u *UserDB) CreateUser(c context.Context, user *entity.User) (err error) {
	if err = DB(c, u.db).Create(user).Error; err != nil {
		return errors.WrapDBError(err, "新增用户失败")
	}
	return err
}

func (u *UserDB) UpdateUser(c context.Context, user *entity.User) (err error) {
	if err = DB(c, u.db).Updates(user).Error; err != nil {
		return errors.WrapDBError(err, "更新用户失败")
	}
	return err
}

func (u *UserDB) GetUserIDsByName(ctx context.Context, username string) (ids []int64, err error) {
	nameQuery := "%" + username + "%"
	if err = DB(ctx, u.db).
		Model(&entity.User{}).
		Select("id").
		Where("username LIKE ?", nameQuery).
		Find(&ids).Error; err != nil {
		err = errors.WrapDBError(err, "获取用户ID失败")
		return
	}
	if ids == nil {
		ids = []int64{}
	}
	return
}

func (u *UserDB) GetUserPage(c context.Context, user *query.User) (userPage []*entity.User, total int64, err error) {
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
		err = errors.WrapDBError(err, "获取用户分页失败")
		return
	}

	// 查询分页数据（可根据需要添加 Limit / Offset 支持）
	if err = queryConditions.Scopes(user.Paginate()).Find(&userPage).Error; err != nil {
		err = errors.WrapDBError(err, "获取用户分页数据失败")
		return
	}

	return
}

func (u *UserDB) GetUserMapByIDs(ctx context.Context, ids []int64) (map[int64]string, error) {
	var users []entity.User
	if err := DB(ctx, u.db).
		Select("id, username").
		Where("id IN ?", ids).
		Find(&users).Error; err != nil {
		return nil, err
	}
	userMap := make(map[int64]string, len(users))
	for _, u := range users {
		userMap[u.ID] = u.Username
	}
	return userMap, nil
}

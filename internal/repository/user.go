package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	myerrors "Art-Design-Backend/pkg/errors"
	"context"
	"errors"
)

type UserRepo struct {
	userDB    *db.UserDB
	userCache *cache.UserCache
	roleCache *cache.RoleCache
}

func NewUserRepo(
	userDB *db.UserDB,
	userCache *cache.UserCache,
	roleCache *cache.RoleCache,
) *UserRepo {
	return &UserRepo{
		userDB:    userDB,
		userCache: userCache,
		roleCache: roleCache,
	}
}

func (u *UserRepo) CheckUserDuplicate(c context.Context, user *entity.User) (err error) {
	return u.userDB.CheckUserDuplicate(c, user)
}

func (u *UserRepo) GetUserById(c context.Context, id int64) (user *entity.User, err error) {
	return u.userDB.GetUserById(c, id)
}

func (u *UserRepo) CreateUser(c context.Context, user *entity.User) (err error) {
	return u.userDB.CreateUser(c, user)
}

func (u *UserRepo) UpdateUser(c context.Context, user *entity.User) (err error) {
	return u.userDB.UpdateUser(c, user)
}

func (u *UserRepo) GetUserPage(c context.Context, user *query.User) (userPage []*entity.User, total int64, err error) {
	return u.userDB.GetUserPage(c, user)
}

func (u *UserRepo) InvalidUserRoleCache(c context.Context, userID int64, originalRoleIds []int64) (err error) {
	var errs []error
	if e := u.roleCache.InvalidRoleUserDepCache(userID, originalRoleIds); e != nil {
		errs = append(errs, e)
	}
	if e := u.userCache.InvalidUserRoleCache(userID); e != nil {
		errs = append(errs, e)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return myerrors.WrapCacheError(err, "用户角色缓存信息失效失败")
}

func (u *UserRepo) GetLoginUserByUsername(c context.Context, username string) (user *entity.User, err error) {
	return u.userDB.GetLoginUserByUsername(c, username)
}

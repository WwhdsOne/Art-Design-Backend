package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	myerrors "Art-Design-Backend/pkg/errors"
	"context"
	"errors"
)

type UserRepo struct {
	*db.UserDB
	*cache.UserCache
	*cache.RoleCache
}

func (u *UserRepo) InvalidUserRoleCache(c context.Context, userID int64, originalRoleIds []int64) (err error) {
	var errs []error
	if e := u.RoleCache.InvalidRoleUserDepCache(userID, originalRoleIds); e != nil {
		errs = append(errs, e)
	}
	if e := u.UserCache.InvalidUserRoleCache(userID); e != nil {
		errs = append(errs, e)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return myerrors.WrapCacheError(err, "用户角色缓存信息失效失败")
}

func (u *UserRepo) GetLoginUserByUsername(c context.Context, username string) (user *entity.User, err error) {
	return u.UserDB.GetLoginUserByUsername(c, username)
}

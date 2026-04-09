package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/cachex"
	"Art-Design-Backend/pkg/constant/rediskey"
	myerrors "Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"fmt"
)

type UserCache struct {
	redis *redisx.RedisWrapper
}

func NewUserCache(redis *redisx.RedisWrapper) *UserCache {
	return &UserCache{
		redis: redis,
	}
}

func (u *UserCache) GetUserRoleList(userID int64) ([]*entity.Role, error) {
	key := fmt.Sprintf("%s%d", rediskey.UserRoleList, userID)
	roleList, err := cachex.GetSlice[*entity.Role](u.redis, key)
	if err != nil || len(roleList) == 0 {
		return nil, myerrors.NewCacheError("获取用户角色信息缓存失败")
	}
	return roleList, nil
}

func (u *UserCache) InvalidUserRoleCache(userID int64) error {
	key := fmt.Sprintf("%s%d", rediskey.UserRoleList, userID)
	if err := u.redis.Del(key); err != nil {
		return myerrors.WrapCacheError(err, "删除用户角色信息缓存失败")
	}
	return nil
}

func (u *UserCache) SetUserRoleList(userID int64, roleList []*entity.Role) error {
	key := fmt.Sprintf("%s%d", rediskey.UserRoleList, userID)
	if err := cachex.Set(u.redis, key, roleList, rediskey.UserRoleListTTL); err != nil {
		return myerrors.WrapCacheError(err, "设置用户角色信息缓存失败")
	}
	return nil
}

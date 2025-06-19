package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"fmt"
	"github.com/bytedance/sonic"
)

type UserCache struct {
	redis *redisx.RedisWrapper
}

func NewUserCache(redis *redisx.RedisWrapper) *UserCache {
	return &UserCache{
		redis: redis,
	}
}

func (u *UserCache) GetUserRoleList(userID int64) (roleList []*entity.Role, err error) {
	key := fmt.Sprintf("%s%d", rediskey.UserRoleList, userID)
	val, err := u.redis.Get(key)
	if err = sonic.Unmarshal([]byte(val), &roleList); err != nil {
		err = errors.NewCacheError("获取用户角色信息缓存失败")
	}
	return
}

func (u *UserCache) InvalidUserRoleCache(userID int64) (err error) {
	userRoleInfoKey := fmt.Sprintf("%s%d", rediskey.UserRoleList, userID)
	if err = u.redis.Del(userRoleInfoKey); err != nil {
		return errors.WrapCacheError(err, "删除用户角色信息缓存失败")
	}
	return
}

func (u *UserCache) SetUserRoleList(userID int64, roleList []*entity.Role) (err error) {
	key := fmt.Sprintf("%s%d", rediskey.UserRoleList, userID)
	val, _ := sonic.Marshal(roleList)
	if err = u.redis.Set(key, string(val), rediskey.UserRoleListTTL); err != nil {
		return errors.WrapCacheError(err, "设置用户角色信息缓存失败")
	}
	return
}

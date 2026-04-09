package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/cachex"
	"Art-Design-Backend/pkg/constant/rediskey"
	myerrors "Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"fmt"
)

type RoleCache struct {
	redis *redisx.RedisWrapper
}

func NewRoleCache(redis *redisx.RedisWrapper) *RoleCache {
	return &RoleCache{
		redis: redis,
	}
}

func (r *RoleCache) InvalidRoleInfoCache(roleID int64) error {
	key := fmt.Sprintf("%s%d", rediskey.RoleUserDependencies, roleID)
	if err := r.redis.DelBySetMembers(key); err != nil {
		return myerrors.WrapCacheError(err, "删除角色信息缓存失败")
	}
	return nil
}

func (r *RoleCache) InvalidRoleUserDepCache(userID int64, originalRoleIDs []int64) error {
	userRoleInfoKey := fmt.Sprintf(rediskey.UserRoleList+"%d", userID)
	for _, roleID := range originalRoleIDs {
		roleUserDepKey := fmt.Sprintf(rediskey.RoleUserDependencies+"%d", roleID)
		if err := r.redis.SRem(roleUserDepKey, userRoleInfoKey); err != nil {
			return myerrors.WrapCacheError(err, "删除用户角色对应关系失败")
		}
	}
	return nil
}

func (r *RoleCache) GetRoleInfo(roleID int64) (*entity.Role, error) {
	key := fmt.Sprintf(rediskey.RoleInfo+"%d", roleID)
	role, err := cachex.GetOrNil[entity.Role](r.redis, key)
	if err != nil {
		return nil, myerrors.WrapCacheError(err, "获取角色信息失败")
	}
	return role, nil
}

func (r *RoleCache) SetRoleInfo(role *entity.Role) error {
	key := fmt.Sprintf(rediskey.RoleInfo+"%d", role.ID)
	if err := cachex.Set(r.redis, key, role, rediskey.RoleInfoTTL); err != nil {
		return myerrors.WrapCacheError(err, "设置角色信息缓存失败")
	}
	return nil
}

func (r *RoleCache) SetRoleUserDep(userID int64, roleID int64) error {
	roleUserDepKey := fmt.Sprintf(rediskey.RoleUserDependencies+"%d", roleID)
	userRoleInfoKey := fmt.Sprintf(rediskey.UserRoleList+"%d", userID)
	if err := r.redis.SAdd(roleUserDepKey, userRoleInfoKey); err != nil {
		return myerrors.WrapCacheError(err, "设置角色用户对应关系失败")
	}
	return nil
}

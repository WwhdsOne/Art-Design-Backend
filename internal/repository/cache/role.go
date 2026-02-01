package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant/rediskey"
	myerrors "Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"fmt"

	"github.com/bytedance/sonic"
)

type RoleCache struct {
	redis *redisx.RedisWrapper
}

func NewRoleCache(redis *redisx.RedisWrapper) *RoleCache {
	return &RoleCache{
		redis: redis,
	}
}

// InvalidRoleInfoCache 删除角色信息缓存
// 同时也删除映射表缓存
func (r *RoleCache) InvalidRoleInfoCache(roleID int64) (err error) {
	// 删除角色信息缓存
	key := fmt.Sprintf("%s%d", rediskey.RoleUserDependencies, roleID)
	// 这里删除了	RoleUserDependencies
	err = r.redis.DelBySetMembers(key)
	if err != nil {
		return myerrors.WrapCacheError(err, "删除角色信息缓存失败")
	}
	return
}

func (r *RoleCache) InvalidRoleUserDepCache(userID int64, originalRoleIDs []int64) (err error) {
	userRoleInfoKey := fmt.Sprintf(rediskey.UserRoleList+"%d", userID)
	for _, roleID := range originalRoleIDs {
		roleUserDepKey := fmt.Sprintf(rediskey.RoleUserDependencies+"%d", roleID)
		if err = r.redis.SRem(roleUserDepKey, userRoleInfoKey); err != nil {
			return myerrors.WrapCacheError(err, "删除用户角色对应关系失败")
		}
	}
	return
}

func (r *RoleCache) GetRoleInfo(roleID int64) (role *entity.Role, err error) {
	key := fmt.Sprintf(rediskey.RoleInfo+"%d", roleID)
	var roleJSON string
	roleJSON, err = r.redis.Get(key)
	if err != nil {
		err = myerrors.WrapCacheError(err, "获取角色信息失败")
		return
	}
	_ = sonic.UnmarshalString(roleJSON, &role)
	return
}

func (r *RoleCache) SetRoleInfo(role *entity.Role) (err error) {
	key := fmt.Sprintf(rediskey.RoleInfo+"%d", role.ID)
	roleJSON, _ := sonic.MarshalString(role)
	err = r.redis.Set(key, roleJSON, rediskey.RoleInfoTTL)
	if err != nil {
		return myerrors.WrapCacheError(err, "设置角色信息缓存失败")
	}
	return
}

func (r *RoleCache) SetRoleUserDep(userID int64, roleID int64) (err error) {
	roleUserDepKey := fmt.Sprintf(rediskey.RoleUserDependencies+"%d", roleID)
	userRoleInfoKey := fmt.Sprintf(rediskey.UserRoleList+"%d", userID)
	err = r.redis.SAdd(roleUserDepKey, userRoleInfoKey)
	if err != nil {
		return myerrors.WrapCacheError(err, "设置角色用户对应关系失败")
	}
	return
}

package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"fmt"
	"github.com/bytedance/sonic"
)

type RoleCache struct {
	Redis *redisx.RedisWrapper
}

func NewRoleCache(redis *redisx.RedisWrapper) *RoleCache {
	return &RoleCache{
		Redis: redis,
	}
}

// InvalidRoleInfoCache 删除角色信息缓存
// 同时也删除映射表缓存
func (r *RoleCache) InvalidRoleInfoCache(roleID int64) (err error) {
	// 删除角色信息缓存
	key := fmt.Sprintf("%s%d", rediskey.RoleInfo, roleID)
	err = r.Redis.DelBySetMembers(key)
	if err != nil {
		return errors.WrapCacheError(err, "删除角色信息缓存失败")
	}
	return
}

func (r *RoleCache) InvalidRoleUserDepCache(userID int64, originalRoleIds []int64) (err error) {
	userRoleInfoKey := fmt.Sprintf(rediskey.UserRoleList+"%d", userID)
	for _, roleID := range originalRoleIds {
		roleUserDepKey := fmt.Sprintf(rediskey.RoleUserDependencies+"%d", roleID)
		if err = r.Redis.SRem(roleUserDepKey, userRoleInfoKey); err != nil {
			return errors.WrapCacheError(err, "删除用户角色对应关系失败")
		}
	}
	return
}

func (r *RoleCache) GetRoleInfo(roleID int64) (role *entity.Role, err error) {
	key := fmt.Sprintf(rediskey.RoleInfo+"%d", roleID)
	var roleJson string
	roleJson, err = r.Redis.Get(key)
	if err != nil {
		err = errors.WrapCacheError(err, "获取角色信息失败")
		return
	}
	_ = sonic.UnmarshalString(roleJson, &role)
	return
}

func (r *RoleCache) SetRoleInfo(role *entity.Role) (err error) {
	key := fmt.Sprintf(rediskey.RoleInfo+"%d", role.ID)
	roleJson, _ := sonic.MarshalString(role)
	err = r.Redis.Set(key, roleJson, rediskey.RoleInfoTTL)
	if err != nil {
		return errors.WrapCacheError(err, "设置角色信息缓存失败")
	}
	return
}

func (r *RoleCache) SetRoleUserDep(roleID int64, userID int64) (err error) {
	roleUserDepKey := fmt.Sprintf(rediskey.RoleUserDependencies+"%d", roleID)
	userRoleInfoKey := fmt.Sprintf(rediskey.UserRoleList+"%d", userID)
	err = r.Redis.SAdd(roleUserDepKey, userRoleInfoKey)
	if err != nil {
		return errors.WrapCacheError(err, "设置角色用户对应关系失败")
	}
	return
}

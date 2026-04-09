package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository/cachex"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"fmt"
	"slices"
	"strings"
)

type MenuCache struct {
	redis *redisx.RedisWrapper
}

func NewMenuCache(redis *redisx.RedisWrapper) *MenuCache {
	return &MenuCache{
		redis: redis,
	}
}

// InvalidateMenuCacheByRoleID 清除与指定角色关联的所有菜单缓存
//
// 缓存设计说明：
// 1. 用户菜单缓存策略：
//   - 每个用户的菜单缓存键格式: "MENU:LIST:ROLE:{roleID1}_{roleID2}_{...}"
//     (例如：用户拥有角色1和2 → "MENU:LIST:ROLE:1_2"
//     用户拥有角色1,2和3 → "MENU:LIST:ROLE:1_2_3")
//
// 2. 反向依赖关系表：
//   - 数据结构：Redis Set
//   - 键格式:   "MENU:ROLE:DEPENDENCIES:{roleID}"
//   - 值内容：  所有包含该roleID的用户菜单缓存键集合
//
// 3. 缓存失效机制：
//   当角色权限变更时：
//   a) 根据 roleID 从 "MENU:ROLE:DEPENDENCIES:{roleID}" 获取所有关联缓存键
//   b) 批量删除这些用户菜单缓存
//   c) 最后清理该角色的依赖记录
func (m *MenuCache) InvalidateMenuCacheByRoleID(roleID int64) error {
	depKey := fmt.Sprintf(rediskey.MenuRoleDependencies+"%d", roleID)
	return m.redis.DelBySetMembers(depKey)
}

func (m *MenuCache) InvalidAllMenuCache() error {
	if err := m.redis.DeleteByPrefix(rediskey.MenuListRole, 100); err != nil {
		return err
	}
	return m.redis.DeleteByPrefix(rediskey.MenuRoleDependencies, 100)
}

func buildMenuCacheKey(roleIDList []int64) string {
	slices.Sort(roleIDList)
	return fmt.Sprintf(rediskey.MenuListRole+"%s", strings.Join(strings.Split(fmt.Sprint(roleIDList), " "), "_"))
}

func (m *MenuCache) GetMenuListByRoleIDListFromCache(roleIDList []int64) ([]*response.Menu, error) {
	key := buildMenuCacheKey(roleIDList)
	return cachex.GetSlice[*response.Menu](m.redis, key)
}

func (m *MenuCache) SetMenuListCache(roleIDList []int64, menuList []*entity.Menu) error {
	key := buildMenuCacheKey(roleIDList)
	if err := cachex.Set(m.redis, key, menuList, rediskey.MenuListRoleTTL); err != nil {
		return errors.WrapCacheError(err, "菜单列表写入缓存失败")
	}
	for _, roleID := range roleIDList {
		depKey := fmt.Sprintf(rediskey.MenuRoleDependencies+"%d", roleID)
		if err := m.redis.SAdd(depKey, key); err != nil {
			return errors.WrapCacheError(err, "设置角色菜单依赖关系失败")
		}
	}
	return nil
}

package cache

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"fmt"
	"sort"
	"strings"

	"github.com/bytedance/sonic"
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
//
//   - 数据结构：Redis Set
//
//   - 键格式:   "MENU:ROLE:DEPENDENCIES:{roleID}"
//
//   - 值内容：  所有包含该roleID的用户菜单缓存键集合
//     (例如："MENU:ROLE:DEPENDENCIES:1" 包含 ["MENU:LIST:ROLE:1_2", "MENU:LIST:ROLE:1_3"])
//
//     3. 缓存失效机制：
//     当角色权限变更时：
//     a) 根据 roleID 从 "MENU:ROLE:DEPENDENCIES:{roleID}" 获取所有关联缓存键
//     b) 批量删除这些用户菜单缓存
//     c) 最后清理该角色的依赖记录
//
// 示例流程：
//   - 用户A(角色1,2) → 缓存键: "MENU:LIST:ROLE:1_2"
//   - 用户B(角色1,3) → 缓存键: "MENU:LIST:ROLE:1_3"
//   - Redis中会建立：
//     "MENU:ROLE:DEPENDENCIES:1" → ["MENU:LIST:ROLE:1_2", "MENU:LIST:ROLE:1_3"]
//     "MENU:ROLE:DEPENDENCIES:2" → ["MENU:LIST:ROLE:1_2"]
//     "MENU:ROLE:DEPENDENCIES:3" → ["MENU:LIST:ROLE:1_3"]
//   - 当角色1权限变更时，自动清除两个用户的菜单缓存，以及角色1的依赖缓存表。
func (m *MenuCache) InvalidateMenuCacheByRoleID(roleID int64) (err error) {
	// 获取记录角色所关联的菜单缓存 key 的依赖集合 key（Set 类型）
	depKey := fmt.Sprintf(rediskey.MenuRoleDependencies+"%d", roleID)

	// 构造删除列表：包括依赖集合本身 和 依赖集合中记录的所有菜单缓存 key
	err = m.redis.DelBySetMembers(depKey)

	return
}

// InvalidAllMenuCache 批量清除所有菜单缓存
func (m *MenuCache) InvalidAllMenuCache() (err error) {
	// 清除所有菜单相关缓存
	if err = m.redis.DeleteByPrefix(rediskey.MenuListRole, 100); err != nil {
		return
	}
	err = m.redis.DeleteByPrefix(rediskey.MenuRoleDependencies, 100)
	return
}

func buildMenuCacheKey(roleIDList []int64) string {
	sort.Slice(roleIDList, func(i, j int) bool {
		return roleIDList[i] < roleIDList[j]
	})
	return fmt.Sprintf(rediskey.MenuListRole+"%s", strings.Join(strings.Split(fmt.Sprint(roleIDList), " "), "_"))
}
func (m *MenuCache) GetMenuListByRoleIDListFromCache(roleIDList []int64) (menu []*response.Menu, err error) {
	key := buildMenuCacheKey(roleIDList)
	val, err := m.redis.Get(key)
	if err != nil {
		return
	}
	err = sonic.Unmarshal([]byte(val), &menu)
	return
}

// SetMenuListCache 缓存菜单列表
func (m *MenuCache) SetMenuListCache(roleIDList []int64, menuList []*entity.Menu) (err error) {
	key := buildMenuCacheKey(roleIDList)
	cacheBytes, err := sonic.Marshal(menuList)
	if err != nil {
		return errors.WrapCacheError(err, "菜单列表序列化失败")
	}

	if err = m.redis.Set(key, string(cacheBytes), rediskey.MenuListRoleTTL); err != nil {
		return errors.WrapCacheError(err, "菜单列表写入缓存失败")
	}
	for _, roleID := range roleIDList {
		depKey := fmt.Sprintf(rediskey.MenuRoleDependencies+"%d", roleID)
		if err = m.redis.SAdd(depKey, key); err != nil {
			return errors.WrapCacheError(err, "设置角色菜单依赖关系失败")
		}
	}
	return
}

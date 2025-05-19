package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/constant/rediskey"
	myerror "Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sort"
	"strings"
	"sync"
)

type MenuService struct {
	MenuRepo      *repository.MenuRepository      // 用户Repo
	RoleRepo      *repository.RoleRepository      // 角色 Repo
	RoleMenusRepo *repository.RoleMenusRepository // 角色菜单关联 Repo
	UserRolesRepo *repository.UserRolesRepository // 用户角色关联 Repo
	Redis         *redisx.RedisWrapper            // redis
	MenuListLocks *sync.Map                       // 根据角色键的菜单列表锁
}

func (m *MenuService) getMenuLock(key string) *sync.RWMutex {
	actual, _ := m.MenuListLocks.LoadOrStore(key, &sync.RWMutex{})
	return actual.(*sync.RWMutex)
}

func (m *MenuService) CreateMenu(c context.Context, menu *request.Menu) (err error) {
	var menuDo entity.Menu
	err = copier.Copy(&menuDo, &menu)
	// 默认父级为顶级菜单
	if menu.ParentID == 0 {
		menuDo.ParentID = -1
	}
	if err != nil {
		zap.L().Error("菜单属性复制失败", zap.Error(err))
		return
	}
	// todo 每个新的菜单都应该跟超级管理员产生关联
	err = m.MenuRepo.CreateMenu(c, &menuDo)
	if err != nil {
		return
	}
	return
}

// GetMenuList 由于每次都是先是获取用户信息，再根据用户角色获取菜单列表
func (m *MenuService) GetMenuList(c context.Context) (res []*response.Menu, err error) {
	// 获取用户角色ID
	userRoleKey := fmt.Sprintf(rediskey.UserRoleList+"%d", authutils.GetUserID(c))
	val, err := m.Redis.Get(userRoleKey)
	// 获取当前用户角色ID
	var roleList []entity.Role
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var roleIDList []int64
			roleIDList, err = m.UserRolesRepo.GetRoleIDListByUserID(c, authutils.GetUserID(c))
			if err != nil {
				zap.L().Error("获取用户角色ID列表失败", zap.Error(err))
				return
			}
			roleList, err = m.RoleRepo.GetRoleListByRoleIDList(c, roleIDList)
			if err != nil {
				zap.L().Error("获取角色列表失败", zap.Error(err))
				return
			}
			zap.L().Debug("获取用户角色对应关系缓存未命中", zap.Int64("userID", authutils.GetUserID(c)))
		} else {
			zap.L().Error("获取用户角色对应关系缓存失败", zap.Error(err))
			err = myerror.NewCacheError("缓存获取失败,请刷新页面")
			return
		}
	}
	_ = sonic.Unmarshal([]byte(val), &roleList)
	roleIds := make([]int64, 0, len(roleList))
	for _, role := range roleList {
		roleIds = append(roleIds, role.ID)
	}
	// 构建缓存 key
	buildMenuCacheKey := func(roleIds []int64) string {
		sort.Slice(roleIds, func(i, j int) bool {
			return roleIds[i] < roleIds[j]
		})
		ids := make([]string, len(roleIds))
		for i, id := range roleIds {
			ids[i] = fmt.Sprintf("%d", id)
		}
		return rediskey.MenuListRole + strings.Join(ids, "_")
	}
	// 提取成局部函数：尝试读取缓存
	tryGetCache := func(key string) bool {
		cacheData, cacheErr := m.Redis.Get(key)
		if cacheErr == nil {
			if unmarshalErr := sonic.Unmarshal([]byte(cacheData), &res); unmarshalErr != nil {
				zap.L().Error("菜单列表缓存解析失败", zap.String("key", key), zap.Error(unmarshalErr))
				return false
			}
			return true
		}
		if errors.Is(cacheErr, redis.Nil) {
			zap.L().Debug("Redis 获取菜单缓存未命中", zap.String("key", key))
		} else {
			zap.L().Error("Redis 获取缓存失败", zap.String("key", key), zap.Error(cacheErr))
			err = cacheErr
		}
		return false
	}

	// 调用局部函数尝试读取缓存
	cacheKey := buildMenuCacheKey(roleIds)
	// 根据角色键获取锁
	lock := m.getMenuLock(cacheKey)

	lock.RLock()
	if tryGetCache(cacheKey) || err != nil {
		lock.RUnlock()
		return
	}
	lock.RUnlock()

	lock.Lock()
	defer lock.Unlock()

	// 再次尝试读取缓存（写前再检查，避免重复构建）
	if tryGetCache(cacheKey) || err != nil {
		return
	}
	// 查询角色关联菜单数据
	menuIDList, err := m.RoleMenusRepo.GetMenuIDListByRoleIDList(c, roleIds)
	if err != nil {
		return
	}
	// 获取数据库数据
	menuList, err := m.MenuRepo.GetMenuListByIDList(c, menuIDList)
	// 先用 map 存储所有菜单，方便查找
	menuMap := make(map[int64]*response.Menu)
	for _, menuDo := range menuList {
		var menuResp response.Menu
		err = copier.Copy(&menuResp, &menuDo)
		if err != nil {
			zap.L().Error("菜单属性复制失败", zap.Error(err))
			return
		}
		// 如果是不是按钮类型，则初始化 AuthList 和 Children
		if menuDo.Type != 3 {
			menuResp.Meta.AuthList = make([]response.AuthMark, 0)
			menuResp.Children = make([]response.Menu, 0)
		}
		menuMap[menuDo.ID] = &menuResp
	}
	// 遍历所有菜单，构建树形结构
	for _, dbMenu := range menuList {
		frontendMenu := menuMap[dbMenu.ID]
		// 跳过按钮类型（按钮的 AuthCode 会挂到父菜单上）
		if dbMenu.Type == 3 { // 按钮类型
			if parent, exists := menuMap[dbMenu.ParentID]; exists {
				// 将按钮的 AuthCode 添加到父菜单的 AuthList
				parent.Meta.AuthList = append(parent.Meta.AuthList, response.AuthMark{
					ID:   dbMenu.ID,
					Name: dbMenu.Title,
					Code: dbMenu.AuthCode,
				})
			}
			continue
		}
		// 如果是顶级菜单，直接添加到结果列表
		if dbMenu.ParentID == -1 {
			res = append(res, frontendMenu)
		} else {
			// 否则挂到父菜单的 Children 下
			if parent, exists := menuMap[dbMenu.ParentID]; exists {
				parent.Children = append(parent.Children, *frontendMenu)
			}
		}
	}
	// 写入缓存
	// 即使空切片也写入缓存，保证缓存的健壮性
	cacheBytes, _ := sonic.Marshal(&res)
	err = m.Redis.Set(cacheKey, string(cacheBytes), rediskey.MenuListRoleTTL)
	if err != nil {
		zap.L().Error("菜单列表写入缓存失败", zap.Error(err))
		return
	}
	// 写入映射表
	for _, rid := range roleIds {
		depKey := fmt.Sprintf(rediskey.MenuRoleDependencies+"%d", rid)
		err = m.Redis.SAdd(depKey, cacheKey)
		if err != nil {
			zap.L().Error("菜单列表写入映射表失败", zap.Error(err))
			return
		}
	}
	return
}

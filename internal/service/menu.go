package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/constant/rediskey"
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
)

type MenuService struct {
	MenuRepo      *repository.MenuRepository      // 用户Repo
	RoleRepo      *repository.RoleRepository      // 角色 Repo
	RoleMenusRepo *repository.RoleMenusRepository // 角色菜单关联 Repo
	UserRolesRepo *repository.UserRolesRepository // 用户角色关联 Repo
	Redis         *redisx.RedisWrapper            // redis
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

func (m *MenuService) GetMenuList(c context.Context) (res []*response.Menu, err error) {
	roleIds := authutils.GetUserRoleIDs(c)
	if len(roleIds) == 0 {
		return
	}
	// 过滤无效角色ID
	validRoleIDList, err := m.RoleRepo.FilterValidRoleIDs(c, roleIds)
	// 过滤并非当前请求用户的角色ID
	// 因为有可能存在用户某个角色已经被删除了，但是jwt内仍然保存着已经被解绑的角色ID
	validRoleIDList, err = m.UserRolesRepo.FilterValidUserRoles(c, roleIds)
	// 构建缓存键函数
	buildMenuCacheKey := func(roleIds []int64) string {
		// 从小到大排序
		sort.Slice(roleIds, func(i, j int) bool {
			return roleIds[i] < roleIds[j]
		})
		ids := make([]string, len(roleIds))
		for i, id := range roleIds {
			ids[i] = fmt.Sprintf("%d", id)
		}
		return rediskey.MenuListRole + strings.Join(ids, "_")
	}
	// 尝试获取缓存
	cacheKey := buildMenuCacheKey(validRoleIDList)
	cacheData, err := m.Redis.Get(cacheKey)
	// 缓存命中，返回
	if err == nil {
		if err = sonic.Unmarshal([]byte(cacheData), &res); err != nil {
			zap.L().Error("菜单列表缓存解析失败", zap.String("key", cacheKey), zap.Error(err))
			// 缓存解析失败，从数据库重新读取，防止数据污染
		} else {
			// 缓存命中且解析成功，返回
			return
		}
	}
	// 缓存未命中，从数据库读取
	if errors.Is(err, redis.Nil) {
		zap.L().Debug("Redis 获取菜单缓存未命中", zap.String("key", cacheKey), zap.Error(err))
	} else {
		zap.L().Error("Redis 获取缓存失败", zap.String("key", cacheKey), zap.Error(err))
		return
	}
	// 查询角色关联菜单数据
	menuIDList, err := m.RoleMenusRepo.GetMenuIDListByRoleIDList(c, validRoleIDList)
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
	cacheBytes, _ := sonic.Marshal(&res)
	err = m.Redis.Set(cacheKey, string(cacheBytes), rediskey.MenuListRoleTTL)
	if err != nil {
		zap.L().Error("菜单列表写入缓存失败", zap.Error(err))
		return
	}
	// 写入映射表
	for _, rid := range validRoleIDList {
		depKey := fmt.Sprintf(rediskey.MenuRoleDependencies+"%d", rid)
		err = m.Redis.SAdd(depKey, cacheKey)
		if err != nil {
			zap.L().Error("菜单列表写入映射表失败", zap.Error(err))
			return
		}
	}
	return
}

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
	"github.com/gin-gonic/gin"
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

func (m *MenuService) invalidAllMenuCache() (err error) {
	// 清除所有菜单相关缓存
	if err = m.Redis.DeleteByPrefix(rediskey.MenuListRole, 100); err != nil {
		return
	}

	if err = m.Redis.DeleteByPrefix(rediskey.MenuRoleDependencies, 100); err != nil {
		return
	}
	return
}

// BuildMenuTree 构建菜单树结构并挂载按钮权限
// 参数 filterHidden 控制是否过滤隐藏菜单（true 过滤，false 不过滤）
func (m *MenuService) BuildMenuTree(menuList []*entity.Menu, filterHidden bool) (res []*response.Menu, err error) {
	menuMap := make(map[int64]*response.Menu)
	for _, menuDo := range menuList {
		if filterHidden && menuDo.IsHide {
			continue
		}
		var menuResp response.Menu
		err = copier.Copy(&menuResp, &menuDo)
		if err != nil {
			zap.L().Error("菜单属性复制失败", zap.Error(err))
			return nil, err
		}
		if menuDo.Type != 3 {
			menuResp.Meta.AuthList = make([]response.MenuAuth, 0)
			menuResp.Children = make([]response.Menu, 0)
		}
		menuMap[menuDo.ID] = &menuResp
	}
	for _, dbMenu := range menuList {
		frontendMenu := menuMap[dbMenu.ID]
		if frontendMenu == nil {
			continue
		}
		if dbMenu.Type == 3 {
			if parent, exists := menuMap[dbMenu.ParentID]; exists {
				parent.Meta.AuthList = append(parent.Meta.AuthList, response.MenuAuth{
					ID:       dbMenu.ID,
					Name:     dbMenu.Title,
					AuthCode: dbMenu.AuthCode,
				})
			}
			continue
		}
		// 顶级菜单直接添加到结果列表
		if dbMenu.ParentID == -1 {
			res = append(res, frontendMenu)
		} else {
			if parent, exists := menuMap[dbMenu.ParentID]; exists {
				parent.Children = append(parent.Children, *frontendMenu)
			}
		}
	}
	return
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
	err = m.MenuRepo.CreateMenu(c, &menuDo)
	if err != nil {
		return
	}
	return
}

func (m *MenuService) CreateMenuAuth(c context.Context, menu *request.MenuAuth) (err error) {
	var menuDo entity.Menu
	err = copier.Copy(&menuDo, &menu)
	if err != nil {
		zap.L().Error("菜单属性复制失败", zap.Error(err))
		return
	}
	err = m.MenuRepo.CreateMenu(c, &menuDo)
	if err != nil {
		return
	}
	return
}

// GetAllMenus 获取全部菜单
func (m *MenuService) GetAllMenus(c context.Context) (res []*response.Menu, err error) {
	menus, err := m.MenuRepo.GetAllMenus(c)
	if err != nil {
		return
	}
	res, err = m.BuildMenuTree(menus, false)
	return
}

// GetMenuList 获取当前用户菜单列表
// 1. 先通过 Redis 缓存获取用户角色列表
// 2. 若未命中，则从数据库获取角色 ID 列表和角色信息，并构造角色菜单缓存键
// 3. 尝试读取菜单缓存数据（使用读写锁），若未命中则从数据库获取菜单信息
// 4. 构建菜单树结构（嵌套子菜单、挂载按钮权限）
// 5. 将菜单结果写入缓存，并更新角色缓存映射表
func (m *MenuService) GetMenuList(c context.Context) (res []*response.Menu, err error) {
	// Step 1. 获取当前用户 ID
	userID := authutils.GetUserID(c)

	// Step 2. 构建用户角色缓存 key
	userRoleKey := fmt.Sprintf(rediskey.UserRoleList+"%d", userID)

	// Step 3. 尝试从 Redis 获取角色信息
	val, err := m.Redis.Get(userRoleKey)

	var roleIds []int64

	// Step 4. 缓存未命中或出错
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Step 4.1 从数据库获取角色 ID
			roleIds, err = m.UserRolesRepo.GetRoleIDListByUserID(c, userID)
			if err != nil {
				zap.L().Error("获取用户角色ID列表失败", zap.Error(err))
				return
			}
			// Step 4.2 验证角色是否启用
			roleIds, _ = m.RoleRepo.FilterValidRoleIDs(c, roleIds)
			zap.L().Debug("获取用户角色对应关系缓存未命中", zap.Int64("userID", userID))
		} else {
			// Step 4.3 Redis 查询出错（非未命中），返回缓存错误提示
			zap.L().Error("获取用户角色对应关系缓存失败", zap.Error(err))
			err = myerror.NewCacheError("缓存获取失败,请刷新页面")
			return
		}
	} else {
		var roleList []entity.Role
		// Step 5. 缓存命中，反序列化用户角色列表
		_ = sonic.Unmarshal([]byte(val), &roleList)

		// Step 6. 提取角色 ID 列表
		roleIds = make([]int64, 0, len(roleList))
		for _, role := range roleList {
			roleIds = append(roleIds, role.ID)
		}
	}

	// Step 7. 构建菜单缓存 key
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

	// Step 8. 定义函数：尝试从缓存获取菜单数据
	tryGetCache := func(key string) bool {
		cacheData, cacheErr := m.Redis.Get(key)
		if cacheErr == nil {
			if unmarshalErr := sonic.Unmarshal([]byte(cacheData), &res); unmarshalErr != nil {
				zap.L().Error("菜单列表缓存解析失败", zap.String("key", key), zap.Error(unmarshalErr))
				return false
			}
			return true // 缓存命中
		}
		if errors.Is(cacheErr, redis.Nil) {
			zap.L().Debug("Redis 获取菜单缓存未命中", zap.String("key", key))
		} else {
			zap.L().Error("Redis 获取缓存失败", zap.String("key", key), zap.Error(cacheErr))
			err = cacheErr
		}
		return false
	}

	// Step 9. 生成菜单缓存 key
	cacheKey := buildMenuCacheKey(roleIds)

	// Step 10. 获取菜单缓存锁（根据角色组合）
	lock := m.getMenuLock(cacheKey)

	// Step 11. 加读锁尝试读取缓存
	lock.RLock()
	if tryGetCache(cacheKey) || err != nil {
		lock.RUnlock()
		return
	}
	lock.RUnlock()

	// Step 12. 加写锁准备写缓存（双检锁）
	lock.Lock()
	defer lock.Unlock()

	// Step 13. 写前再次检查缓存（双检锁模式）
	if tryGetCache(cacheKey) || err != nil {
		return
	}

	// Step 14. 查询当前角色对应的菜单 ID 列表
	menuIDList, err := m.RoleMenusRepo.GetMenuIDListByRoleIDList(c, roleIds)
	if err != nil {
		return
	}

	// Step 15. 根据菜单 ID 列表获取菜单实体
	menuList, err := m.MenuRepo.GetMenuListByIDList(c, menuIDList)
	if err != nil {
		return
	}

	// Step 16. 构建菜单树结构
	res, err = m.BuildMenuTree(menuList, true)

	// Step 17. 将结果写入 Redis 缓存
	cacheBytes, _ := sonic.Marshal(&res)
	err = m.Redis.Set(cacheKey, string(cacheBytes), rediskey.MenuListRoleTTL)
	if err != nil {
		zap.L().Error("菜单列表写入缓存失败", zap.Error(err))
		return
	}

	// Step 18. 写入菜单依赖映射表（用于后续清缓存）
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

func (m *MenuService) DeleteMenu(c *gin.Context, id int64) (err error) {
	// Step 1. 获取所有需要删除的菜单 ID（包括子菜单、按钮）
	var allMenuIDs []int64
	err = m.collectMenuIDTree(c, id, &allMenuIDs)
	if err != nil {
		zap.L().Error("递归收集子菜单失败", zap.Error(err))
		return
	}

	if len(allMenuIDs) == 0 {
		return
	}

	// Step 2. 删除菜单表中记录
	err = m.MenuRepo.DeleteMenuByIDList(c, allMenuIDs)
	if err != nil {
		zap.L().Error("删除菜单失败", zap.Error(err))
		return
	}

	// Step 3. 删除角色-菜单映射关系
	err = m.RoleMenusRepo.DeleteByMenuIDs(c, allMenuIDs)
	if err != nil {
		zap.L().Error("删除角色菜单关系失败", zap.Error(err))
		return
	}

	// Step 4. 清除所有菜单相关缓存
	go func() {
		if err = m.invalidAllMenuCache(); err != nil {
			zap.L().Error("删除菜单权限缓存失败", zap.Error(err))
		}
	}()

	return
}

// collectMenuIDTree 递归收集某个菜单及其所有子菜单 ID
func (m *MenuService) collectMenuIDTree(c context.Context, parentID int64, result *[]int64) (err error) {
	*result = append(*result, parentID)

	children, err := m.MenuRepo.GetChildMenuIDsByParentID(c, parentID)
	if err != nil {
		return err
	}

	for _, childID := range children {
		err = m.collectMenuIDTree(c, childID, result)
		if err != nil {
			return err
		}
	}
	return
}

func (m *MenuService) UpdateMenuAuth(c *gin.Context, r *request.MenuAuth) (err error) {
	var menu entity.Menu
	err = copier.Copy(&menu, r)
	if err != nil {
		zap.L().Error("权限参数复制失败", zap.Error(err))
		return
	}
	err = m.MenuRepo.CheckMenuDuplicate(&menu)
	if err != nil {
		zap.L().Error("菜单权限重复", zap.Error(err))
		return
	}
	err = m.MenuRepo.UpdateMenuAuth(c, &menu)
	if err != nil {
		zap.L().Error("更新菜单权限失败", zap.Error(err))
		return
	}
	go func() {
		if err = m.invalidAllMenuCache(); err != nil {
			zap.L().Error("删除菜单权限缓存失败", zap.Error(err))
		}
	}()
	return
}

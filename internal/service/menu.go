package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/authutils"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type MenuService struct {
	MenuRepo *repository.MenuRepo // 菜单Repo
	RoleRepo *repository.RoleRepo // 角色Repo
}

// buildMenuTree 构建菜单树结构并挂载按钮权限
// 参数 filterHidden 控制是否过滤隐藏菜单（true 过滤，false 不过滤）
func (m *MenuService) buildMenuTree(menuList []*entity.Menu) (res []*response.Menu, err error) {
	menuMap := make(map[int64]*response.Menu)
	for _, menuDo := range menuList {
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
					AuthCode: *dbMenu.AuthCode,
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
		zap.L().Error("创建菜单失败", zap.Error(err))
		return
	}
	return
}

func (m *MenuService) UpdateMenu(c context.Context, r *request.Menu) (err error) {
	var menuDo entity.Menu
	err = copier.Copy(&menuDo, &r)
	if err != nil {
		zap.L().Error("菜单属性复制失败", zap.Error(err))
		return
	}
	err = m.MenuRepo.CheckMenuDuplicate(c, &menuDo)
	if err != nil {
		zap.L().Error("更新菜单时,菜单名称重复", zap.Error(err))
		return
	}
	err = m.MenuRepo.UpdateMenu(c, &menuDo)
	if err != nil {
		zap.L().Error("更新菜单失败", zap.Error(err))
		return
	}
	go func() {
		if err := m.MenuRepo.InvalidAllMenuCache(); err != nil {
			zap.L().Error("更新菜单时,删除菜单缓存失败", zap.Error(err))
		}
	}()
	return
}

// CreateMenuAuth 创建菜单权限
func (m *MenuService) CreateMenuAuth(c context.Context, menu *request.MenuAuth) (err error) {
	var menuDo entity.Menu
	err = copier.Copy(&menuDo, &menu)
	if err != nil {
		zap.L().Error("创建菜单权限时,菜单属性复制失败", zap.Error(err))
		return
	}
	err = m.MenuRepo.CreateMenu(c, &menuDo)
	if err != nil {
		zap.L().Error("创建菜单权限失败", zap.Error(err))
		return
	}
	return
}

// UpdateMenuAuth 更新菜单权限
func (m *MenuService) UpdateMenuAuth(c *gin.Context, r *request.MenuAuth) (err error) {
	var menu entity.Menu
	err = copier.Copy(&menu, r)
	if err != nil {
		zap.L().Error("权限参数复制失败", zap.Error(err))
		return
	}
	err = m.MenuRepo.CheckMenuDuplicate(c, &menu)
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
		if err := m.MenuRepo.InvalidAllMenuCache(); err != nil {
			zap.L().Error("删除菜单权限缓存失败", zap.Error(err))
		}
	}()
	return
}

// GetAllMenus 获取全部菜单
func (m *MenuService) GetAllMenus(c context.Context) (res []*response.Menu, err error) {
	menus, err := m.MenuRepo.GetAllMenus(c)
	if err != nil {
		return
	}
	res, err = m.buildMenuTree(menus)
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
	// Step 2. 获取当前用户角色列表
	var roleIDList []int64
	roleIDList, err = m.RoleRepo.GetRoleIDListByUserID(c, userID)

	// Step 3. 尝试从缓存获取当前角色菜单列表
	res, err = m.MenuRepo.GetMenuListByRoleIDListFromCache(roleIDList)
	if err == nil {
		return
	}

	// Step 4. 降级,根据菜单 ID 列表获取菜单实体
	menuList, err := m.MenuRepo.GetMenuListByRoleIDList(c, roleIDList)
	if err != nil {
		return
	}

	// Step 5. 构建菜单树结构
	res, err = m.buildMenuTree(menuList)

	// Step 6. 将结果写入 Redis 缓存
	go func(roleIDList []int64, menuList []*entity.Menu) {
		if err := m.MenuRepo.SetMenuListCache(c, roleIDList, menuList); err != nil {
			zap.L().Warn("写入菜单缓存失败", zap.Error(err))
		}
	}(roleIDList, menuList)
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
	if err = m.MenuRepo.DeleteMenuByIDList(c, allMenuIDs); err != nil {
		zap.L().Error("删除菜单失败", zap.Error(err))
		return
	}

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

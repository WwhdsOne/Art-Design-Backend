package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"context"
	"go.uber.org/zap"
)

type MenuRepo struct {
	menuCache  *cache.MenuCache
	menuDB     *db.MenuDB
	roleMenuDB *db.RoleMenusDB
}

func NewMenuRepo(
	menuDB *db.MenuDB,
	menuCache *cache.MenuCache,
	roleMenuDB *db.RoleMenusDB,
) *MenuRepo {
	return &MenuRepo{
		menuDB:     menuDB,
		menuCache:  menuCache,
		roleMenuDB: roleMenuDB,
	}
}

func (m *MenuRepo) GetReducedMenuList(c context.Context) (menuList []*entity.Menu, err error) {
	return m.menuDB.GetReducedMenuList(c)
}

func (m *MenuRepo) InvalidateMenuCacheByRoleID(roleID int64) (err error) {
	return m.menuCache.InvalidateMenuCacheByRoleID(roleID)
}

func (m *MenuRepo) InvalidAllMenuCache(c context.Context) (err error) {
	return m.menuCache.InvalidAllMenuCache()
}

func (m *MenuRepo) GetMenuIDListByRoleID(c context.Context, roleID int64) (menuIDList []int64, err error) {
	return m.menuDB.GetMenuIDListByRoleID(c, roleID)
}

func (m *MenuRepo) CreateMenu(c context.Context, menu *entity.Menu) (err error) {
	return m.menuDB.CreateMenu(c, menu)
}

func (m *MenuRepo) UpdateMenu(c context.Context, menu *entity.Menu) (err error) {
	return m.menuDB.UpdateMenu(c, menu)
}

func (m *MenuRepo) CheckMenuDuplicate(c context.Context, menu *entity.Menu) (err error) {
	return m.menuDB.CheckMenuDuplicate(c, menu)
}

func (m *MenuRepo) UpdateMenuAuth(c context.Context, menu *entity.Menu) (err error) {
	return m.menuDB.UpdateMenuAuth(c, menu)
}

func (m *MenuRepo) GetAllMenus(c context.Context) (res []*entity.Menu, err error) {
	return m.menuDB.GetAllMenus(c)
}
func (m *MenuRepo) GetMenuListByIDList(c context.Context, menuIDList []int64) (menuList []*entity.Menu, err error) {
	return m.menuDB.GetMenuListByIDList(c, menuIDList)
}

func (m *MenuRepo) DeleteMenuByIDList(c context.Context, menuIDList []int64) (err error) {
	if err = m.menuDB.DeleteMenuByIDList(c, menuIDList); err != nil {
		return
	}
	if err = m.roleMenuDB.DeleteRoleMenuRelationByMenuIDs(c, menuIDList); err != nil {
		return
	}
	go func() {
		// 影响权限的必须报错，不能用警告
		if err = m.menuCache.InvalidAllMenuCache(); err != nil {
			zap.L().Error("删除菜单权限缓存失败", zap.Error(err))
		}
	}()
	return
}

func (m *MenuRepo) GetChildMenuIDsByParentID(c context.Context, parentID int64) (childrenIDs []int64, err error) {
	return m.menuDB.GetChildMenuIDsByParentID(c, parentID)
}

func (m *MenuRepo) GetMenuListByRoleIDListFromCache(c context.Context, roleIDList []int64) (menuList []*response.Menu, err error) {
	return m.menuCache.GetMenuListByRoleIDList(roleIDList)
}

func (m *MenuRepo) GetMenuIDListByRoleIDList(c context.Context, roleIDList []int64) (menuIDList []int64, err error) {
	return m.menuDB.GetMenuIDListByRoleIDList(c, roleIDList)
}

func (m *MenuRepo) GetMenuListByRoleIDList(c context.Context, roleIDList []int64) (menuList []*entity.Menu, err error) {
	menuIDList, err := m.menuDB.GetMenuIDListByRoleIDList(c, roleIDList)
	if err != nil {
		return
	}
	return m.menuDB.GetMenuListByIDList(c, menuIDList)
}

func (m *MenuRepo) SetMenuListCache(c context.Context, roleIDList []int64, menuList []*entity.Menu) (err error) {
	return m.menuCache.SetMenuListCache(roleIDList, menuList)
}

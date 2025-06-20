package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"context"
	"go.uber.org/zap"
)

type MenuRepo struct {
	*cache.MenuCache
	*db.MenuDB
	*db.RoleMenusDB
}

func (m *MenuRepo) DeleteMenuByIDList(c context.Context, menuIDList []int64) (err error) {
	if err = m.MenuDB.DeleteMenuByIDList(c, menuIDList); err != nil {
		return
	}
	if err = m.RoleMenusDB.DeleteRoleMenuRelationByMenuIDs(c, menuIDList); err != nil {
		return
	}
	go func() {
		// 影响权限的必须报错，不能用警告
		if err := m.MenuCache.InvalidAllMenuCache(); err != nil {
			zap.L().Error("删除菜单权限缓存失败", zap.Error(err))
		}
	}()
	return
}

func (m *MenuRepo) GetMenuListByRoleIDList(c context.Context, roleIDList []int64) (menuList []*entity.Menu, err error) {
	menuIDList, err := m.MenuDB.GetMenuIDListByRoleIDList(c, roleIDList)
	if err != nil {
		return
	}
	return m.MenuDB.GetMenuListByIDList(c, menuIDList)
}

func (m *MenuRepo) SetMenuListCache(c context.Context, roleIDList []int64, menuList []*entity.Menu) (err error) {
	return m.MenuCache.SetMenuListCache(roleIDList, menuList)
}

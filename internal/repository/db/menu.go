package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type MenuDB struct {
	db *gorm.DB // 用户表数据库连接
}

func NewMenuDB(db *gorm.DB) *MenuDB {
	return &MenuDB{
		db: db,
	}
}

func (m *MenuDB) CheckMenuDuplicate(c context.Context, menu *entity.Menu) (err error) {
	var result struct {
		NameExists  bool
		PathExists  bool
		TitleExists bool
	}

	excludeID := ""
	if menu.ID != 0 {
		excludeID = fmt.Sprintf("AND id != %d", menu.ID)
	}

	var query strings.Builder
	args := make([]interface{}, 0)
	conditions := make([]string, 0)

	if menu.Name != nil && *menu.Name != "" {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM \"menu\" WHERE name = ? "+excludeID+") AS name_exists")
		args = append(args, *menu.Name)
	}

	if menu.Path != nil && *menu.Path != "" {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM \"menu\" WHERE path = ? "+excludeID+") AS path_exists")
		args = append(args, *menu.Path)
	}

	if menu.Title != "" {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM \"menu\" WHERE title = ? "+excludeID+") AS title_exists")
		args = append(args, menu.Title)
	}

	if len(conditions) == 0 {
		return nil // 没有需要查重的字段
	}

	query.WriteString("SELECT ")
	query.WriteString(strings.Join(conditions, ", "))

	if err = DB(c, m.db).Raw(query.String(), args...).Scan(&result).Error; err != nil {
		return err
	}

	switch {
	case result.NameExists:
		err = errors.NewDBError("组件名称重复")
	case result.PathExists:
		err = errors.NewDBError("路由地址重复")
	case result.TitleExists:
		err = errors.NewDBError("菜单名称重复")
	}
	return
}

func (m *MenuDB) GetAllMenus(c context.Context) (res []*entity.Menu, err error) {
	if err = DB(c, m.db).Find(&res).Error; err != nil {
		return nil, errors.WrapDBError(err, "获取所有菜单失败")
	}
	return
}

func (m *MenuDB) GetMenuListByIDList(c context.Context, menuIDList []int64) (menuList []*entity.Menu, err error) {
	if err = DB(c, m.db).
		Where("id IN ?", menuIDList).
		Find(&menuList).Error; err != nil {
		err = errors.WrapDBError(err, "获取菜单失败")
		return
	}
	return
}

func (m *MenuDB) GetReducedMenuList(c context.Context) (menuList []*entity.Menu, err error) {
	if err = DB(c, m.db).
		Select("id", "title", "parent_id", "type").
		Find(&menuList).Error; err != nil {
		err = errors.WrapDBError(err, "获取菜单失败")
		return
	}
	return
}

func (m *MenuDB) CreateMenu(c context.Context, menu *entity.Menu) (err error) {
	if err = DB(c, m.db).Create(menu).Error; err != nil {
		zap.L().Error("创建菜单失败", zap.Error(err))
		return errors.NewDBError("创建菜单失败")
	}
	return
}

func (m *MenuDB) DeleteMenuByIDList(c context.Context, menuIDList []int64) (err error) {
	if err = DB(c, m.db).Where("id IN ?", menuIDList).Delete(&entity.Menu{}).Error; err != nil {
		zap.L().Error("删除菜单失败", zap.Error(err))
		return errors.NewDBError("删除菜单失败")
	}
	return
}

func (m *MenuDB) GetChildMenuIDsByParentID(c context.Context, parentID int64) (childrenIDs []int64, err error) {
	if err = DB(c, m.db).Model(&entity.Menu{}).
		Where("parent_id = ?", parentID).
		Pluck("id", &childrenIDs).Error; err != nil {
		return nil, errors.WrapDBError(err, "获取子菜单失败")
	}
	return
}

func (m *MenuDB) UpdateMenuAuth(c context.Context, menu *entity.Menu) (err error) {
	if err = DB(c, m.db).Where("id = ?", menu.ID).Updates(menu).Error; err != nil {
		return errors.NewDBError("更新菜单失败")
	}
	return
}

func (m *MenuDB) UpdateMenu(c context.Context, menu *entity.Menu) (err error) {
	if err = DB(c, m.db).Where("id = ?", menu.ID).Updates(menu).Error; err != nil {
		return errors.NewDBError("更新菜单失败")
	}
	return
}

func (m *MenuDB) GetMenuIDListByRoleIDList(c context.Context, roleIDList []int64) (menuIDList []int64, err error) {
	if err = DB(c, m.db).
		Model(&entity.RoleMenus{}).
		Where("role_id IN ?", roleIDList).
		Pluck("menu_id", &menuIDList).Error; err != nil {
		err = errors.WrapDBError(err, "获取角色菜单关联信息失败")
		return
	}
	return
}

func (m *MenuDB) GetMenuIDListByRoleID(c context.Context, roleID int64) (menuIDList []int64, err error) {
	if err = DB(c, m.db).
		Model(&entity.RoleMenus{}).
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDList).Error; err != nil {
		err = errors.WrapDBError(err, "获取角色菜单关联信息失败")
		return
	}
	return
}

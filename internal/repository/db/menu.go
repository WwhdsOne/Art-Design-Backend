package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/dupcheck"
	"Art-Design-Backend/pkg/errors"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MenuDB struct {
	db *gorm.DB // 用户表数据库连接
}

func NewMenuDB(db *gorm.DB) *MenuDB {
	return &MenuDB{
		db: db,
	}
}

func (m *MenuDB) CheckMenuDuplicate(c context.Context, menu *entity.Menu) error {
	name, nameIsNil := "", menu.Name == nil
	if !nameIsNil {
		name = *menu.Name
	}
	path, pathIsNil := "", menu.Path == nil
	if !pathIsNil {
		path = *menu.Path
	}
	return dupcheck.Check(DB(c, m.db), "\"menu\"", menu.ID, []dupcheck.Field{
		{Column: "name", ErrMsg: "组件名称重复", Value: name, IsNil: nameIsNil},
		{Column: "path", ErrMsg: "路由地址重复", Value: path, IsNil: pathIsNil},
		{Column: "title", ErrMsg: "菜单名称重复", Value: menu.Title},
	})
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
	/*
		⚠️ 这里必须使用 map[string]interface{} 进行更新，而不能直接用 Updates(struct)

		原因：
		GORM 在使用 Updates(struct) 时，只会更新“非零值”字段。
		而 Go 中 bool 的零值是 false，
		如果 IsInMainContainer = false，
		GORM 会认为这是零值，从而忽略该字段，不会写入数据库。

		例如：
			IsInMainContainer = false
		将不会出现在 UPDATE SQL 里。

		因此这里改为 map 更新方式：
			map 更新会无条件更新指定字段（包括 false、0、"" 等零值）

		这是生产环境推荐写法，避免 bool/int/string 零值被静默忽略。
	*/

	updateData := map[string]any{
		"name":                 menu.Name,
		"type":                 menu.Type,
		"path":                 menu.Path,
		"component":            menu.Component,
		"parent_id":            menu.ParentID,
		"title":                menu.Title,
		"icon":                 menu.Icon,
		"show_badge":           menu.ShowBadge,
		"show_text_badge":      menu.ShowTextBadge,
		"is_hide":              menu.IsHide,
		"is_hide_tab":          menu.IsHideTab,
		"link":                 menu.Link,
		"is_iframe":            menu.IsIframe,
		"keep_alive":           menu.KeepAlive,
		"auth_code":            menu.AuthCode,
		"is_in_main_container": menu.IsInMainContainer,
		"sort":                 menu.Sort,
	}

	if err = DB(c, m.db).
		Model(&entity.Menu{}).
		Where("id = ?", menu.ID).
		Updates(updateData).Error; err != nil {
		return errors.NewDBError("更新菜单失败")
	}

	return nil
}

package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"

	"gorm.io/gorm"
)

type RoleMenusDB struct {
	db *gorm.DB // 用户表数据库连接
}

func NewRoleMenusDB(db *gorm.DB) *RoleMenusDB {
	return &RoleMenusDB{
		db: db,
	}
}

func (r *RoleMenusDB) GetMenuIDListByRoleIDList(c context.Context, roleIDList []int64) (menuIDList []int64, err error) {
	if err = DB(c, r.db).
		Model(&entity.RoleMenus{}).
		Where("role_id IN ?", roleIDList).
		Pluck("menu_id", &menuIDList).Error; err != nil {
		err = errors.WrapDBError(err, "获取角色菜单关联信息失败")
		return
	}
	return
}

func (r *RoleMenusDB) GetMenuIDListByRoleID(c context.Context, roleID int64) (menuIDList []int64, err error) {
	if err = DB(c, r.db).
		Model(&entity.RoleMenus{}).
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDList).Error; err != nil {
		err = errors.WrapDBError(err, "获取角色菜单关联信息失败")
		return
	}
	return
}

func (r *RoleMenusDB) DeleteMenuRelationsByRoleID(c context.Context, roleID int64) (err error) {
	if err = DB(c, r.db).
		Where("role_id = ?", roleID).
		Delete(&entity.RoleMenus{}).Error; err != nil {
		return errors.WrapDBError(err, "删除角色菜单关联失败")
	}
	return
}

// CreateRoleMenus 创建角色菜单关联
func (r *RoleMenusDB) CreateRoleMenus(c context.Context, roleID int64, menuIDList []int64) (err error) {
	roleMenus := make([]entity.RoleMenus, 0, len(menuIDList))
	for _, menuID := range menuIDList {
		roleMenus = append(roleMenus, entity.RoleMenus{
			RoleID: roleID,
			MenuID: menuID,
		})
	}
	if err = DB(c, r.db).Create(&roleMenus).Error; err != nil {

		return errors.WrapDBError(err, "创建角色菜单关联失败")
	}
	return
}

// DeleteRoleMenuRelationByMenuIDs 删除角色菜单关联
func (r *RoleMenusDB) DeleteRoleMenuRelationByMenuIDs(c context.Context, menuIDList []int64) (err error) {
	if err = DB(c, r.db).
		Where("menu_id IN ?", menuIDList).
		Delete(&entity.RoleMenus{}).Error; err != nil {
		return errors.NewDBError("删除角色菜单关联失败")
	}
	return
}

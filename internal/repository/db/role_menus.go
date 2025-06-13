package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleMenusRepository struct {
	db *gorm.DB // 用户表数据库连接
}

func NewRoleMenusRepository(db *gorm.DB) *RoleMenusRepository {
	return &RoleMenusRepository{
		db: db,
	}
}

func (r *RoleMenusRepository) GetMenuIDListByRoleIDList(c context.Context, roleIDList []int64) (menuIDList []int64, err error) {
	if err = DB(c, r.db).
		Model(&entity.RoleMenus{}).
		Where("role_id IN ?", roleIDList).
		Pluck("menu_id", &menuIDList).Error; err != nil {
		zap.L().Error("获取角色菜单关联信息失败", zap.Error(err))
		err = errors.NewDBError("获取角色菜单关联信息失败")
		return
	}
	return
}

func (r *RoleMenusRepository) GetMenuIDListByRoleID(c context.Context, roleID int64) (menuIDList []int64, err error) {
	if err = DB(c, r.db).
		Model(&entity.RoleMenus{}).
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDList).Error; err != nil {
		zap.L().Error("获取角色菜单关联信息失败", zap.Error(err))
		err = errors.NewDBError("获取角色菜单关联信息失败")
		return
	}
	return
}

func (r *RoleMenusRepository) DeleteByRoleID(c context.Context, roleID int64) (err error) {
	if err = DB(c, r.db).
		Where("role_id = ?", roleID).
		Delete(&entity.RoleMenus{}).Error; err != nil {
		zap.L().Error("删除角色菜单关联失败", zap.Error(err))
		return errors.NewDBError("删除角色菜单关联失败")
	}
	return
}

// CreateRoleMenus 创建角色菜单关联
// 由于创建只会在删除后进行，所以创建函数不调整缓存
func (r *RoleMenusRepository) CreateRoleMenus(c context.Context, roleID int64, menuIDList []int64) (err error) {
	roleMenus := make([]entity.RoleMenus, 0, len(menuIDList))
	for _, menuID := range menuIDList {
		roleMenus = append(roleMenus, entity.RoleMenus{
			RoleID: roleID,
			MenuID: menuID,
		})
	}
	if err = DB(c, r.db).Create(&roleMenus).Error; err != nil {
		zap.L().Error("创建角色菜单关联失败", zap.Error(err))
		return errors.NewDBError("创建角色菜单关联失败")
	}
	return
}

// DeleteByMenuIDs 删除角色菜单关联
func (r *RoleMenusRepository) DeleteByMenuIDs(c *gin.Context, menuIDList []int64) (err error) {
	if err = DB(c, r.db).
		Where("menu_id IN ?", menuIDList).
		Delete(&entity.RoleMenus{}).Error; err != nil {
		zap.L().Error("删除角色菜单关联失败", zap.Error(err))
		return errors.NewDBError("删除角色菜单关联失败")
	}
	return
}

package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/errors"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB // 用户表数据库连接
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{
		db: db.Table(constant.MenuTableName),
	}
}

func (m *MenuRepository) GetMenuList(c context.Context) (menuList []*entity.Menu, err error) {
	if err = DB(c, m.db).Find(menuList).Error; err != nil {
		zap.L().Error("获取菜单失败", zap.Error(err))
		return nil, errors.NewGormError("获取菜单失败")
	}
	return
}

func (m *MenuRepository) GetMenuListByRoleIDList(c context.Context, roleIdList []int64) (menuList []*entity.Menu, err error) {
	if err = DB(c, m.db).
		Select("menu.*").
		Joins("JOIN role_menus ON role_menus.menu_id = menu.id").
		Where("role_menus.role_id IN ?", roleIdList).
		Find(&menuList).Error; err != nil {
		zap.L().Error("获取菜单失败", zap.Error(err))
		return nil, errors.NewGormError("获取菜单失败")
	}
	return menuList, nil
}

func (m *MenuRepository) CreateMenu(c context.Context, menu *entity.Menu) (err error) {
	if err = DB(c, m.db).Create(menu).Error; err != nil {
		zap.L().Error("创建菜单失败", zap.Error(err))
		return errors.NewGormError("创建菜单失败")
	}
	return
}

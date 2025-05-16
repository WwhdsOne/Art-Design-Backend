package repository

import (
	"Art-Design-Backend/internal/model/entity"
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
		db: db,
	}
}

func (m *MenuRepository) GetMenuListByIDList(c context.Context, menuIDList []int64) (menuList []*entity.Menu, err error) {
	if err = DB(c, m.db).
		Where("id IN ?", menuIDList).
		Find(&menuList).Error; err != nil {
		zap.L().Error("获取菜单失败", zap.Error(err))
		err = errors.NewDBError("获取菜单失败")
		return
	}
	return
}

func (m *MenuRepository) CreateMenu(c context.Context, menu *entity.Menu) (err error) {
	if err = DB(c, m.db).Create(menu).Error; err != nil {
		zap.L().Error("创建菜单失败", zap.Error(err))
		return errors.NewDBError("创建菜单失败")
	}
	return
}

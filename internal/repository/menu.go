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

func (m *MenuRepository) GetReducedMenuList(c context.Context) (menuList []*entity.Menu, err error) {
	if err = DB(c, m.db).
		Select("id", "title", "parent_id", "type").
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

func (m *MenuRepository) DeleteMenuByIDList(c context.Context, menuIDList []int64) (err error) {
	if err = DB(c, m.db).Where("id IN ?", menuIDList).Delete(&entity.Menu{}).Error; err != nil {
		zap.L().Error("删除菜单失败", zap.Error(err))
		return errors.NewDBError("删除菜单失败")
	}
	return
}

func (m *MenuRepository) GetChildMenuIDsByParentID(ctx context.Context, parentID int64) (childrenIDs []int64, err error) {
	if err = DB(ctx, m.db).Model(&entity.Menu{}).
		Where("parent_id = ?", parentID).
		Pluck("id", &childrenIDs).Error; err != nil {
		zap.L().Error("获取子菜单失败", zap.Error(err))
		return nil, errors.NewDBError("获取子菜单失败")
	}
	return
}

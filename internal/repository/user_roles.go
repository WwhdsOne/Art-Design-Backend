package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRolesRepository struct {
	db *gorm.DB // 用户表数据库连接
}

func NewUserRolesRepository(db *gorm.DB) *UserRolesRepository {
	return &UserRolesRepository{
		db: db,
	}
}

// GetRoleIDListByUserID 根据用户ID获取角色ID列表
func (u *UserRolesRepository) GetRoleIDListByUserID(c context.Context, userID int64) (roleIDs []int64, err error) {
	if err = DB(c, u.db).
		Model(&entity.UserRoles{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error; err != nil {
		zap.L().Error("查询用户角色列表失败", zap.Error(err))
		err = errors.NewDBError("查询用户角色列表失败")
		return
	}
	return
}

// DeleteRolesFromUserByUserID 删除用户原有角色关联
func (u *UserRolesRepository) DeleteRolesFromUserByUserID(c context.Context, userID int64) (err error) {
	if err = DB(c, u.db).
		Where("user_id = ?", userID).
		Delete(&entity.UserRoles{}).Error; err != nil {
		zap.L().Error("删除原有关联失败", zap.Int64("userID", userID), zap.Error(err))
		return errors.NewDBError("删除原有关联失败")
	}
	return
}

// AddRolesToUser 添加新的角色关联
func (u *UserRolesRepository) AddRolesToUser(c context.Context, userID int64, roleIDList []int64) (err error) {
	if len(roleIDList) == 0 {
		return
	}
	userRoleList := make([]entity.UserRoles, 0, len(roleIDList))
	for _, roleID := range roleIDList {
		userRoleList = append(userRoleList, entity.UserRoles{
			UserID: userID,
			RoleID: roleID,
		})
	}
	if err = DB(c, u.db).Create(&userRoleList).Error; err != nil {
		zap.L().Error("创建新的关联失败", zap.Int64("userID", userID), zap.Error(err))
		return errors.NewDBError("创建新的关联失败")
	}
	return
}

func (u *UserRolesRepository) GetReducedRoleList(ctx context.Context) (roleList []*entity.Role, err error) {
	if err = DB(ctx, u.db).
		Select("id", "name").
		Where("status = 1").
		Find(&roleList).Error; err != nil {
		err = errors.NewDBError("获取精简角色列表失败")
		return
	}
	return
}

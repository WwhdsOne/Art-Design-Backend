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

// AssignRoleToUser 分配角色给用户
func (u *UserRolesRepository) AssignRoleToUser(c context.Context, userID int64, roleIDList []int64) (err error) {
	// 删除原有关联
	if err = DB(c, u.db).
		Where("user_id = ?", userID).
		Delete(&entity.UserRoles{}).Error; err != nil {
		zap.L().Error("删除原有关联失败")
		err = errors.NewDBError("删除原有关联失败")
		return
	}

	// 创建新的关联
	if len(roleIDList) > 0 {
		userRoleList := make([]entity.UserRoles, 0, len(roleIDList))
		for _, roleID := range roleIDList {
			userRoleList = append(userRoleList, entity.UserRoles{
				UserID: userID,
				RoleID: roleID,
			})
		}
		if err = DB(c, u.db).
			Create(userRoleList).Error; err != nil {
			zap.L().Error("创建新的关联失败")
			err = errors.NewDBError("创建新的关联失败")
			return
		}
	}

	return
}

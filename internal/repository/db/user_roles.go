package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"

	"gorm.io/gorm"
)

type UserRolesDB struct {
	db *gorm.DB // 用户表数据库连接
}

func NewUserRolesDB(db *gorm.DB) *UserRolesDB {
	return &UserRolesDB{
		db: db,
	}
}

// DeleteUserRoleRelationsByUserID 根据用户ID删除用户角色关联
func (u *UserRolesDB) DeleteUserRoleRelationsByUserID(c context.Context, userID int64) (err error) {
	if err = DB(c, u.db).
		Where("user_id = ?", userID).
		Delete(&entity.UserRoles{}).Error; err != nil {
		return errors.WrapDBError(err, "删除原有关联失败")
	}
	return
}

// GetRoleIDListByUserID 根据用户ID获取角色ID列表
func (u *UserRolesDB) GetRoleIDListByUserID(c context.Context, userID int64) (roleIDs []int64, err error) {
	if err = DB(c, u.db).
		Model(&entity.UserRoles{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error; err != nil {
		err = errors.WrapDBError(err, "查询用户角色列表失败")
		return
	}
	return
}

// AddRolesToUser 添加新的角色关联
func (u *UserRolesDB) AddRolesToUser(c context.Context, userID int64, roleIDList []int64) (err error) {
	userRoleList := make([]*entity.UserRoles, 0, len(roleIDList))
	for _, roleID := range roleIDList {
		userRoleList = append(userRoleList, &entity.UserRoles{
			UserID: userID,
			RoleID: roleID,
		})
	}
	if err = DB(c, u.db).Create(&userRoleList).Error; err != nil {
		return errors.WrapDBError(err, "创建新的关联失败")
	}
	return
}

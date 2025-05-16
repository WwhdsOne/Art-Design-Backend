package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRolesRepository struct {
	db    *gorm.DB             // 用户表数据库连接
	redis *redisx.RedisWrapper // redis缓存
}

func NewUserRolesRepository(db *gorm.DB, redis *redisx.RedisWrapper) *UserRolesRepository {
	return &UserRolesRepository{
		db:    db,
		redis: redis,
	}
}

// FilterValidUserRoles 获取用户实际拥有的有效角色ID列表
func (u *UserRolesRepository) FilterValidUserRoles(c context.Context, roleIDs []int64) (validRoleIDs []int64, err error) {
	if err = DB(c, u.db).
		Model(&entity.UserRoles{}).
		Where("user_id = ? AND role_id IN ?", authutils.GetUserID(c), roleIDs).
		Pluck("role_id", &validRoleIDs).Error; err != nil {
		zap.L().Error("查询用户有效角色失败", zap.Error(err))
		err = errors.NewDBError("查询用户有效角色失败")
		return
	}
	return
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

	// 删除原有角色缓存
	cacheKey := fmt.Sprintf(rediskey.UserRoleList+"%d", userID)
	if delErr := u.redis.Del(cacheKey); delErr != nil {
		zap.L().Warn("删除用户角色缓存失败", zap.Int64("userID", userID))
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
			Table("user_roles").
			Create(userRoleList).Error; err != nil {
			zap.L().Error("创建新的关联失败")
			err = errors.NewDBError("创建新的关联失败")
			return
		}
	}

	return
}

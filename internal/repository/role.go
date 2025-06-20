package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"context"
	"go.uber.org/zap"
)

type RoleRepo struct {
	*db.RoleDB
	*db.RoleMenusDB
	*cache.RoleCache
	*cache.UserCache
	*db.UserRolesDB
}

func (r *RoleRepo) GetRoleIDListByUserID(c context.Context, userID int64) (roleIDs []int64, err error) {
	list, err := r.UserCache.GetUserRoleList(userID)
	if err != nil {
		zap.L().Warn("获取用户角色列表缓存错误", zap.Error(err))
	} else {
		for _, role := range list {
			roleIDs = append(roleIDs, role.ID)
		}
		return
	}
	roleIDs, err = r.UserRolesDB.GetRoleIDListByUserID(c, userID)
	if err != nil {
		return
	}
	return
}

func (r *RoleRepo) GetRoleListByUserID(c context.Context, userID int64) (roleList []*entity.Role, err error) {
	// 1. 尝试从 Redis 获取缓存
	roleList, err = r.UserCache.GetUserRoleList(userID)
	if err == nil {
		return
	}

	// 2. 查询角色ID列表
	roleIDList, err := r.GetRoleIDListByUserID(c, userID)
	if err != nil {
		return
	}

	// 3. 根据ID列表去数据库或缓存查询数据
	for _, roleID := range roleIDList {
		var role *entity.Role
		role, err = r.RoleCache.GetRoleInfo(roleID)
		if err == nil {
			roleList = append(roleList, role)
			continue
		}
		// 3.2 从数据库读取
		role, err = r.RoleDB.GetEnableRoleByID(c, roleID)
		if err != nil {
			return
		}
		// 3.3 角色信息异步写入 Redis
		go func(role *entity.Role) {
			if err := r.RoleCache.SetRoleInfo(role); err != nil {
				zap.L().Warn("角色缓存写入失败", zap.Error(err))
			}
		}(role)
		roleList = append(roleList, role)
	}

	// 4. 用户角色对应关系写入 Redis 缓存
	go func(userID int64, roleList []*entity.Role) {
		if err := r.UserCache.SetUserRoleList(userID, roleList); err != nil {
			zap.L().Warn("用户角色对应关系写入缓存失败", zap.Int64("userID", userID), zap.Error(err))
		}
		// 5. 写入 Redis 映射表（每个角色映射该用户缓存key）
		for _, roleID := range roleIDList {
			if err := r.RoleCache.SetRoleUserDep(userID, roleID); err != nil {
				zap.L().Warn("用户角色对应关系写入映射表失败", zap.Error(err))
			}
		}
	}(userID, roleList)
	return
}

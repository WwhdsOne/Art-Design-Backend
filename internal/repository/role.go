package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"context"
	"go.uber.org/zap"
)

type RoleRepo struct {
	roleDB      *db.RoleDB
	roleMenusDB *db.RoleMenusDB
	roleCache   *cache.RoleCache
	userCache   *cache.UserCache
	userRoleDB  *db.UserRolesDB
}

func NewRoleRepo(
	roleDB *db.RoleDB,
	roleMenusDB *db.RoleMenusDB,
	roleCache *cache.RoleCache,
	userCache *cache.UserCache,
	userRoleDB *db.UserRolesDB,
) *RoleRepo {
	return &RoleRepo{
		roleDB:      roleDB,
		roleMenusDB: roleMenusDB,
		roleCache:   roleCache,
		userCache:   userCache,
		userRoleDB:  userRoleDB,
	}
}

func (r *RoleRepo) CheckRoleDuplicate(c context.Context, role *entity.Role) (err error) {
	return r.roleDB.CheckRoleDuplicate(c, role)
}

func (r *RoleRepo) GetRoleIDListByUserID(c context.Context, userID int64) (roleIDs []int64, err error) {
	list, err := r.userCache.GetUserRoleList(userID)
	if err != nil {
		zap.L().Warn("获取用户角色列表缓存错误", zap.Error(err))
	} else {
		for _, role := range list {
			roleIDs = append(roleIDs, role.ID)
		}
		return
	}
	roleIDs, err = r.userRoleDB.GetRoleIDListByUserID(c, userID)
	if err != nil {
		return
	}
	return
}

func (r *RoleRepo) GetRoleListByUserID(c context.Context, userID int64) (roleList []*entity.Role, err error) {
	// 1. 尝试从 Redis 获取缓存
	roleList, err = r.userCache.GetUserRoleList(userID)
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
		role, err = r.roleCache.GetRoleInfo(roleID)
		if err == nil {
			roleList = append(roleList, role)
			continue
		}
		// 3.2 从数据库读取
		role, err = r.roleDB.GetEnableRoleByID(c, roleID)
		if err != nil {
			return
		}
		// 3.3 角色信息异步写入 Redis
		go func(role *entity.Role) {
			if err := r.roleCache.SetRoleInfo(role); err != nil {
				zap.L().Warn("角色缓存写入失败", zap.Error(err))
			}
		}(role)
		roleList = append(roleList, role)
	}

	// 4. 用户角色对应关系写入 Redis 缓存
	go func(userID int64, roleList []*entity.Role) {
		if err := r.userCache.SetUserRoleList(userID, roleList); err != nil {
			zap.L().Warn("用户角色对应关系写入缓存失败", zap.Int64("userID", userID), zap.Error(err))
		}
		// 5. 写入 Redis 映射表（每个角色映射该用户缓存key）
		for _, roleID := range roleIDList {
			if err := r.roleCache.SetRoleUserDep(userID, roleID); err != nil {
				zap.L().Warn("用户角色对应关系写入映射表失败", zap.Error(err))
			}
		}
	}(userID, roleList)
	return
}

func (r *RoleRepo) GetReducedRoleList(ctx context.Context) (roleList []*entity.Role, err error) {
	return r.roleDB.GetReducedRoleList(ctx)
}

func (r *RoleRepo) DeleteUserRoleRelationsByUserID(c context.Context, userID int64) (err error) {
	return r.userRoleDB.DeleteUserRoleRelationsByUserID(c, userID)
}

func (r *RoleRepo) AddRolesToUser(c context.Context, userID int64, roleIDList []int64) (err error) {
	return r.userRoleDB.AddRolesToUser(c, userID, roleIDList)
}

func (r *RoleRepo) CreateRole(c context.Context, role *entity.Role) (err error) {
	return r.roleDB.CreateRole(c, role)
}

func (r *RoleRepo) GetEnableRoleByID(c context.Context, roleID int64) (role *entity.Role, err error) {
	return r.roleDB.GetEnableRoleByID(c, roleID)
}

func (r *RoleRepo) GetRolePage(c context.Context, role *query.Role) (rolePage []*entity.Role, total int64, err error) {
	return r.roleDB.GetRolePage(c, role)
}

func (r *RoleRepo) UpdateRole(c context.Context, role *entity.Role) (err error) {
	return r.roleDB.UpdateRole(c, role)
}

func (r *RoleRepo) DeleteRoleByID(c context.Context, roleID int64) (err error) {
	return r.roleDB.DeleteRoleByID(c, roleID)
}

func (r *RoleRepo) InvalidRoleInfoCache(c context.Context, roleID int64) (err error) {
	return r.roleCache.InvalidRoleInfoCache(roleID)
}

func (r *RoleRepo) DeleteMenuRelationsByRoleID(c context.Context, roleID int64) (err error) {
	return r.roleMenusDB.DeleteMenuRelationsByRoleID(c, roleID)
}

func (r *RoleRepo) CreateRoleMenus(c context.Context, roleID int64, menuIDList []int64) (err error) {
	return r.roleMenusDB.CreateRoleMenus(c, roleID, menuIDList)
}

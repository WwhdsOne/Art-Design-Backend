package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleMenusRepository struct {
	db    *gorm.DB             // 用户表数据库连接
	redis *redisx.RedisWrapper // redis缓存
}

func NewRoleMenusRepository(db *gorm.DB, redis *redisx.RedisWrapper) *RoleMenusRepository {
	return &RoleMenusRepository{
		db:    db,
		redis: redis,
	}
}

// invalidateMenuCacheByRoleID 清除与某个角色绑定的所有菜单缓存
// 缓存策略说明：
//
// 我们为每个用户构建菜单时，会将其所拥有的所有角色联合作为缓存 key（如：user:menus:{roleIDs}）进行缓存；
// 为了能在角色权限变更时正确删除这些用户级缓存，需要维护一个角色 → 菜单缓存 key 的「反向依赖关系表」：
//   - 该依赖表结构使用 Redis 的 Set 实现，Key: menu:deps:role:{roleID}，Value: 所有使用了该 roleID 的菜单缓存 key（即 user:menus:{roleIDs}）
//   - 当角色的菜单更新时，我们只需要根据 menu:deps:role:{roleID} 中的成员，批量删除相关缓存即可，避免缓存污染
//
// 示例：
//   假设用户 A 有角色 [1, 2]，其缓存 key 为 user:menus:[1,2]
//   那么 menu:deps:role:1 和 menu:deps:role:2 中都应包含 user:menus:[1,2]
//   当角色 1 或 2 被修改菜单后，都可以触发删除该缓存，保障数据一致性

func (r *RoleMenusRepository) invalidateMenuCacheByRoleID(roleID int64) (err error) {
	// 获取记录角色所关联的菜单缓存 key 的依赖集合 key（Set 类型）
	depKey := fmt.Sprintf(rediskey.MenuRoleDependencies+"%d", roleID)

	// 查询该角色所依赖的所有用户菜单缓存 key
	cacheKeys := r.redis.SMembers(depKey)

	// 构造删除列表：包括依赖集合本身 和 依赖集合中记录的所有菜单缓存 key
	delKeys := make([]string, 0, len(cacheKeys)+1)
	delKeys = append(delKeys, depKey)       // 删除依赖表（防止过期错误等影响新依赖写入）
	delKeys = append(delKeys, cacheKeys...) // 删除所有受影响的用户菜单缓存

	// 批量删除（使用 Redis 管道提高性能）
	if err = r.redis.PipelineDel(delKeys); err != nil {
		zap.L().Error("角色菜单缓存批量删除失败", zap.Int64("roleID", roleID), zap.Error(err))
		return
	}

	return
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
	if err = r.invalidateMenuCacheByRoleID(roleID); err != nil {
		zap.L().Error("删除角色菜单关联缓存失败", zap.Error(err))
		return errors.NewDBError("删除角色菜单关联缓存失败")
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

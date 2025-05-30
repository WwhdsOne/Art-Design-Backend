package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
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

// InvalidateMenuCacheByRoleID 清除与指定角色关联的所有菜单缓存
//
// 缓存设计说明：
// 1. 用户菜单缓存策略：
//   - 每个用户的菜单缓存键格式: "MENU:LIST:ROLE:{roleID1}_{roleID2}_{...}"
//     (例如：用户拥有角色1和2 → "MENU:LIST:ROLE:1_2"
//     用户拥有角色1,2和3 → "MENU:LIST:ROLE:1_2_3")
//
// 2. 反向依赖关系表：
//
//   - 数据结构：Redis Set
//
//   - 键格式:   "MENU:ROLE:DEPENDENCIES:{roleID}"
//
//   - 值内容：  所有包含该roleID的用户菜单缓存键集合
//     (例如："MENU:ROLE:DEPENDENCIES:1" 包含 ["MENU:LIST:ROLE:1_2", "MENU:LIST:ROLE:1_3"])
//
//     3. 缓存失效机制：
//     当角色权限变更时：
//     a) 根据 roleID 从 "MENU:ROLE:DEPENDENCIES:{roleID}" 获取所有关联缓存键
//     b) 批量删除这些用户菜单缓存
//     c) 最后清理该角色的依赖记录
//
// 示例流程：
//   - 用户A(角色1,2) → 缓存键: "MENU:LIST:ROLE:1_2"
//   - 用户B(角色1,3) → 缓存键: "MENU:LIST:ROLE:1_3"
//   - Redis中会建立：
//     "MENU:ROLE:DEPENDENCIES:1" → ["MENU:LIST:ROLE:1_2", "MENU:LIST:ROLE:1_3"]
//     "MENU:ROLE:DEPENDENCIES:2" → ["MENU:LIST:ROLE:1_2"]
//     "MENU:ROLE:DEPENDENCIES:3" → ["MENU:LIST:ROLE:1_3"]
//   - 当角色1权限变更时，自动清除两个用户的菜单缓存，以及角色1的依赖缓存表。
func (r *RoleMenusRepository) InvalidateMenuCacheByRoleID(roleID int64) (err error) {
	// 获取记录角色所关联的菜单缓存 key 的依赖集合 key（Set 类型）
	depKey := fmt.Sprintf(rediskey.MenuRoleDependencies+"%d", roleID)

	// 构造删除列表：包括依赖集合本身 和 依赖集合中记录的所有菜单缓存 key
	err = r.redis.DelBySetMembers(depKey)

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

// DeleteByMenuIDs 删除角色菜单关联
func (r *RoleMenusRepository) DeleteByMenuIDs(c *gin.Context, menuIDList []int64) (err error) {
	if err = DB(c, r.db).
		Where("menu_id IN ?", menuIDList).
		Delete(&entity.RoleMenus{}).Error; err != nil {
		zap.L().Error("删除角色菜单关联失败", zap.Error(err))
		return errors.NewDBError("删除角色菜单关联失败")
	}
	return
}

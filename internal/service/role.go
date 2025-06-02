package service

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"strconv"
)

type RoleService struct {
	RoleRepo      *repository.RoleRepository         // 用户Repo
	MenuRepo      *repository.MenuRepository         // 菜单Repo
	RoleMenusRepo *repository.RoleMenusRepository    // 角色菜单Repo
	GormTX        *repository.GormTransactionManager // 事务
	Redis         *redisx.RedisWrapper               // redis
}

// invalidateMenuCacheByRoleID 清除与指定角色关联的所有菜单缓存
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
func (r *RoleService) invalidateMenuCacheByRoleID(roleID int64) (err error) {
	// 获取记录角色所关联的菜单缓存 key 的依赖集合 key（Set 类型）
	depKey := rediskey.MenuRoleDependencies + strconv.FormatInt(roleID, 10)

	// 构造删除列表：包括依赖集合本身 和 依赖集合中记录的所有菜单缓存 key
	err = r.Redis.DelBySetMembers(depKey)

	return
}

// invalidRoleInfoCache 删除角色信息缓存
// 同时也删除映射表缓存
func (r *RoleService) invalidRoleInfoCache(roleID int64) (err error) {
	// 删除角色信息缓存
	key := rediskey.RoleInfo + strconv.FormatInt(roleID, 10)
	err = r.Redis.DelBySetMembers(key)
	if err != nil {
		err = errors.NewCacheError("删除角色信息缓存失败")
		return
	}
	return
}

func (r *RoleService) CreateRole(c context.Context, role *request.Role) (err error) {
	var roleDo entity.Role
	err = copier.Copy(&roleDo, &role)
	if err != nil {
		zap.L().Error("角色属性复制失败", zap.Error(err))
		return
	}
	err = r.RoleRepo.CheckRoleDuplicate(c, &roleDo)
	if err != nil {
		return
	}
	err = r.RoleRepo.CreateRole(c, &roleDo)
	if err != nil {
		return
	}
	return
}

func (r *RoleService) GetRolePage(c *gin.Context, roleQuery *query.Role) (rolePageRes *base.PaginationResp[response.Role], err error) {
	rolePage, total, err := r.RoleRepo.GetRolePage(c, roleQuery)
	if err != nil {
		return
	}
	roleList := make([]response.Role, 0, len(rolePage))
	for _, role := range rolePage {
		var roleResp response.Role
		if err = copier.Copy(&roleResp, &role); err != nil {
			zap.L().Error("复制属性失败")
			return
		}
		roleList = append(roleList, roleResp)
	}
	rolePageRes = base.BuildPageResp[response.Role](roleList, total, roleQuery.PaginationReq)
	return
}

func (r *RoleService) UpdateRole(c *gin.Context, roleReq *request.Role) (err error) {
	var roleDo entity.Role
	if err = copier.Copy(&roleDo, &roleReq); err != nil {
		zap.L().Error("复制属性失败")
		return
	}
	err = r.RoleRepo.CheckRoleDuplicate(c, &roleDo)
	if err != nil {
		return
	}
	if err = r.RoleRepo.UpdateRole(c, &roleDo); err != nil {
		return
	}
	go func() {
		if cacheErr := r.invalidRoleInfoCache(roleDo.ID); cacheErr != nil {
			zap.L().Error("用户信息缓存删除失败（需补偿）", zap.Int64("roleID", roleDo.ID), zap.Error(cacheErr))
		}
	}()
	return
}

func (r *RoleService) DeleteRoleByID(c *gin.Context, roleID int64) (err error) {
	err = r.GormTX.Transaction(c, func(ctx context.Context) error {
		if err = r.RoleRepo.DeleteRoleByID(ctx, roleID); err != nil {
			return err
		}
		if err = r.RoleMenusRepo.DeleteByRoleID(ctx, roleID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	go func() {
		if cacheErr := r.invalidateMenuCacheByRoleID(roleID); cacheErr != nil {
			zap.L().Error("缓存删除失败（需补偿）", zap.Int64("roleID", roleID), zap.Error(cacheErr))
		}
	}()
	return
}

func (r *RoleService) GetRoleMenuBinding(c *gin.Context, roleID int64) (res *response.RoleMenuBinding, err error) {
	res = &response.RoleMenuBinding{}
	var simpleMenuList []*response.SimpleMenu
	menuList, err := r.MenuRepo.GetReducedMenuList(c)
	if err != nil {
		return
	}
	hasMenuIDList, err := r.RoleMenusRepo.GetMenuIDListByRoleID(c, roleID)
	// 先用 map 存储所有菜单，方便查找
	menuMap := make(map[int64]*response.SimpleMenu)
	for _, menuDo := range menuList {
		var menuResp response.SimpleMenu
		err = copier.Copy(&menuResp, &menuDo)
		if err != nil {
			zap.L().Error("菜单属性复制失败", zap.Error(err))
			return
		}
		if menuDo.Type != 3 {
			menuResp.Children = make([]*response.SimpleMenu, 0)
		}
		menuMap[menuDo.ID] = &menuResp
	}
	// 遍历所有菜单，构建树形结构
	for _, dbMenu := range menuList {
		frontendMenu := menuMap[dbMenu.ID]
		// 如果是顶级菜单，直接添加到结果列表
		if dbMenu.ParentID == -1 {
			simpleMenuList = append(simpleMenuList, frontendMenu)
		} else {
			if parent, exists := menuMap[dbMenu.ParentID]; exists {
				parent.Children = append(parent.Children, frontendMenu)
			}
		}
	}
	res.Menus = simpleMenuList
	res.HasMenuIDs = hasMenuIDList
	return
}

func (r *RoleService) UpdateRoleMenuBinding(c *gin.Context, req *request.RoleMenuBinding) (err error) {
	err = r.GormTX.Transaction(c, func(ctx context.Context) error {
		if err = r.RoleMenusRepo.DeleteByRoleID(ctx, int64(req.RoleId)); err != nil {
			return err
		}
		if err = r.RoleMenusRepo.CreateRoleMenus(ctx, int64(req.RoleId), req.MenuIds); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}

	// 缓存清理移出事务
	go func() {
		if cacheErr := r.invalidateMenuCacheByRoleID(int64(req.RoleId)); cacheErr != nil {
			zap.L().Error("缓存删除失败（需补偿）", zap.Int64("roleID", int64(req.RoleId)), zap.Error(cacheErr))
		}
	}()

	return
}

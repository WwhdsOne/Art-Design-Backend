package service

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type RoleService struct {
	RoleRepo      *repository.RoleRepository         // 用户Repo
	MenuRepo      *repository.MenuRepository         // 菜单Repo
	RoleMenusRepo *repository.RoleMenusRepository    // 角色菜单Repo
	GormTX        *repository.GormTransactionManager // 事务
	Redis         *redisx.RedisWrapper               // redis
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
	err = r.GormTX.Transaction(c, func(ctx context.Context) (err error) {
		if err = r.RoleMenusRepo.DeleteByRoleID(ctx, int64(req.RoleId)); err != nil {
			return
		}
		if err = r.RoleMenusRepo.CreateRoleMenus(ctx, int64(req.RoleId), req.MenuIds); err != nil {
			return
		}
		return
	})
	return
}

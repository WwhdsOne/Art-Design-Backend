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

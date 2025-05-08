package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type RoleService struct {
	RoleRepo *repository.RoleRepository // 用户Repo
	Redis    *redisx.RedisWrapper       // redis
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

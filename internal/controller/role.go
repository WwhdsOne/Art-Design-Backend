package controller

import (
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService *service.RoleService // 创建一个MenuService实例
}

func NewRoleController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.RoleService) *RoleController {
	menuCtrl := &RoleController{
		roleService: service,
	}
	r := engine.Group("/api").Group("/role").Use(middleware.AuthMiddleware())
	{
		r.POST("/create", menuCtrl.createRole)
	}
	return menuCtrl
}

func (r *RoleController) createRole(c *gin.Context) {
	var role request.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.Error(err)
		return
	}
	err := r.roleService.CreateRole(c, &role)
	if err != nil {
		c.Error(err)
		return
	}
	response.OkWithMessage("添加成功", c)
}

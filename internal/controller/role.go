package controller

import (
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"
	"Art-Design-Backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService *service.RoleService // 创建一个MenuService实例
}

func NewRoleController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.RoleService) *RoleController {
	menuCtrl := &RoleController{
		roleService: service,
	}
	r := engine.Group("/api").Group("/role")
	r.Use(middleware.AuthMiddleware())
	{
		r.POST("/create", menuCtrl.createRole)
		r.POST("/update", menuCtrl.updateRole)
		r.POST("/page", menuCtrl.gerRolePage)
		r.POST("/delete/:id", menuCtrl.deleteRole)
	}
	return menuCtrl
}

func (r *RoleController) createRole(c *gin.Context) {
	var role request.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		_ = c.Error(err)
		c.Set(gin.BindKey, role)
		return
	}
	err := r.roleService.CreateRole(c, &role)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("添加成功", c)
}

func (r *RoleController) gerRolePage(c *gin.Context) {
	var roleQuery query.Role
	if err := c.ShouldBindJSON(&roleQuery); err != nil {
		_ = c.Error(err)
		return
	}
	rolePageData, err := r.roleService.GetRolePage(c, &roleQuery)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(rolePageData, c)
}

func (r *RoleController) updateRole(c *gin.Context) {
	var role request.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		_ = c.Error(err)
		c.Set(gin.BindKey, role)
		return
	}
	err := r.roleService.UpdateRole(c, &role)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("修改成功", c)
}

func (r *RoleController) deleteRole(c *gin.Context) {
	roleID, err := utils.ParseID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = r.roleService.DeleteRoleByID(c, roleID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("删除成功", c)
}

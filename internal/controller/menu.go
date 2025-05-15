package controller

import (
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"
	"github.com/gin-gonic/gin"
)

type MenuController struct {
	menuService *service.MenuService // 创建一个MenuService实例
}

func NewMenuController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.MenuService) *MenuController {
	menuCtrl := &MenuController{
		menuService: service,
	}
	r := engine.Group("/api").Group("/menu")
	r.Use(middleware.AuthMiddleware())
	{
		// 私有路由组（需要 JWT 认证）
		r.GET("/list", menuCtrl.getMenuList)
		r.POST("/create", menuCtrl.createMenu)
	}
	return menuCtrl
}

func (m *MenuController) getMenuList(c *gin.Context) {
	res, err := m.menuService.GetMenuList(c)
	if err != nil {
		c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

func (m *MenuController) createMenu(c *gin.Context) {
	var menu request.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.Error(err)
		return
	}
	err := m.menuService.CreateMenu(c, &menu)
	if err != nil {
		c.Error(err)
		return
	}
	result.OkWithMessage("添加成功", c)
}

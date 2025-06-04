package controller

import (
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"
	"Art-Design-Backend/pkg/utils"
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
		r.GET("/all", menuCtrl.getAllMenus)
		r.POST("/createMenu", menuCtrl.createMenu)
		r.POST("/updateMenu", menuCtrl.updateMenu)
		r.POST("/createAuth", menuCtrl.createAuth)
		r.POST("/updateAuth", menuCtrl.updateAuth)
		r.POST("/delete/:id", menuCtrl.deleteMenu)

	}
	return menuCtrl
}

func (m *MenuController) getMenuList(c *gin.Context) {
	res, err := m.menuService.GetMenuList(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

func (m *MenuController) createMenu(c *gin.Context) {
	var menu request.Menu
	if err := c.ShouldBindBodyWithJSON(&menu); err != nil {
		_ = c.Error(err)
		return
	}
	err := m.menuService.CreateMenu(c, &menu)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("添加成功", c)
}

func (m *MenuController) createAuth(c *gin.Context) {
	var menu request.MenuAuth
	if err := c.ShouldBindBodyWithJSON(&menu); err != nil {
		_ = c.Error(err)
		return
	}
	err := m.menuService.CreateMenuAuth(c, &menu)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("添加成功", c)
}

func (m *MenuController) updateMenu(c *gin.Context) {
	var menu request.Menu
	if err := c.ShouldBindBodyWithJSON(&menu); err != nil {
		_ = c.Error(err)
		return
	}
	err := m.menuService.UpdateMenu(c, &menu)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("更新成功", c)
}

func (m *MenuController) updateAuth(c *gin.Context) {
	var menu request.MenuAuth
	if err := c.ShouldBindBodyWithJSON(&menu); err != nil {
		_ = c.Error(err)
		return
	}
	err := m.menuService.UpdateMenuAuth(c, &menu)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("更新成功", c)
}

func (m *MenuController) deleteMenu(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = m.menuService.DeleteMenu(c, id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("删除成功", c)
}

func (m *MenuController) getAllMenus(c *gin.Context) {
	res, err := m.menuService.GetAllMenus(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(res, c)
}

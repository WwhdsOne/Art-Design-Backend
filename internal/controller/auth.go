package controller

import (
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService // 创建一个AuthService实例
}

func NewAuthController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.AuthService) *AuthController {
	authCtrl := &AuthController{
		authService: service,
	}
	r := engine.Group("/api").Group("/auth")
	{
		// 私有路由组（需要 JWT 认证）
		r.POST("/logout", middleware.AuthMiddleware(), authCtrl.logout)
	}
	{
		// 公共路由组（无需认证）
		r.GET("/refreshToken/:id", authCtrl.refreshToken)
		r.POST("/register", authCtrl.register)
		r.POST("/login", authCtrl.login)
	}
	return authCtrl
}

// login 处理用户登录请求
func (a *AuthController) login(c *gin.Context) {
	var loginReq request.Login
	err := c.ShouldBindJSON(&loginReq)
	// 如果绑定过程中出现错误，返回错误响应并结束函数执行
	if err != nil {
		c.Error(err)
		c.Set(gin.BindKey, loginReq)
		return
	}
	// 调用service.Login函数尝试验证用户登录信息
	token, err := a.authService.Login(c, &loginReq)
	if err != nil {
		c.Error(err)
		return
	}
	// 返回生成的token
	response.OkWithData(token, c)
}

// logout 处理用户注销请求
func (a *AuthController) logout(c *gin.Context) {
	// 调用 jwt 包中的 LogoutToken 函数注销 token
	err := a.authService.LogoutToken(c)
	if err != nil {
		response.FailWithMessage("注销失败", c)
		return
	}
	response.OkWithMessage("注销成功", c)
}

// refreshToken 处理用户刷新 token 请求
func (a *AuthController) refreshToken(c *gin.Context) {
	token, err := a.authService.RefreshToken(c)
	if err != nil {
		c.Error(err)
		return
	}
	response.OkWithData(token, c)
}

// register 处理用户注册请求
func (a *AuthController) register(c *gin.Context) {
	var userReq request.RegisterUser
	err := c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		c.Set(gin.BindKey, userReq)
		return
	}
	err = a.authService.Register(c, &userReq)
	if err != nil {
		c.Error(err)
		return
	}
	response.OkWithMessage("注册成功", c)
}

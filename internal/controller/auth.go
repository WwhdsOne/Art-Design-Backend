package controller

import (
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService // 创建一个AuthService实例
}

func NewAuthController(engine *gin.Engine, mws *middleware.Middlewares, svc *service.AuthService) *AuthController {
	authCtrl := &AuthController{
		authService: svc,
	}
	r := engine.Group("/api").Group("/auth")
	{
		// 私有路由组（需要 JWT 认证）
		r.POST("/logout", mws.AuthMiddleware(), authCtrl.logout)
	}
	{
		// 公共路由组（无需认证）
		r.POST("/register", authCtrl.register)
		r.POST("/login", authCtrl.login)
	}
	return authCtrl
}

// login 处理用户登录请求
func (a *AuthController) login(c *gin.Context) {
	var loginReq request.Login
	// 如果绑定过程中出现错误，返回错误响应并结束函数执行
	if err := c.ShouldBindBodyWithJSON(&loginReq); err != nil {
		_ = c.Error(err)
		return
	}
	// 调用service.Login函数尝试验证用户登录信息
	token, err := a.authService.Login(c, &loginReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	// 返回生成的token
	result.OkWithData(token, c)
}

// logout 处理用户注销请求
func (a *AuthController) logout(c *gin.Context) {
	// 调用 jwt 包中的 LogoutToken 函数注销 token
	err := a.authService.LogoutToken(c)
	if err != nil {
		result.FailWithMessage("注销失败", c)
		return
	}
	result.OkWithMessage("注销成功", c)
}

// register 处理用户注册请求
func (a *AuthController) register(c *gin.Context) {
	var userReq request.RegisterUser
	if err := c.ShouldBindBodyWithJSON(&userReq); err != nil {
		_ = c.Error(err)
		return
	}
	if err := a.authService.Register(c, &userReq); err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("注册成功", c)
}

package controller

import (
	"Art-Design-Backend/model/request"
	"Art-Design-Backend/pkg/auth"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/response"
	"Art-Design-Backend/service"
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
		r.POST("/refreshToken", middleware.AuthMiddleware(), authCtrl.refreshToken)
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
	// 定义一个Login结构体变量来绑定请求中的JSON数据
	var loginReq request.Login
	// 使用c.ShouldBindJSON尝试将请求中的JSON数据绑定到loginReq变量
	err := c.ShouldBindJSON(&loginReq)
	// 如果绑定过程中出现错误，返回错误响应并结束函数执行
	if err != nil {
		c.Error(err)
		c.Set(gin.BindKey, loginReq)
		return
	}
	// 调用service.Login函数尝试验证用户登录信息
	token, err := a.authService.Login(c, loginReq)
	if err != nil {
		c.Error(err)
		return
	}
	// 返回生成的token
	response.OkWithData(token, c)
}

// logout 处理用户注销请求
func (a *AuthController) logout(c *gin.Context) {
	// 从请求头中获取 token
	token := auth.GetToken(c)
	// 调用 jwt 包中的 LogoutToken 函数注销 token
	err := a.authService.LogoutToken(token)
	if err != nil {
		response.FailWithMessage("注销失败", c)
		return
	}
	response.OkWithMessage("注销成功", c)
}

// refreshToken 处理用户刷新 token 请求
func (a *AuthController) refreshToken(c *gin.Context) {
	// 从上下文中获取用户 ID
	id := auth.GetUserID(c)
	//  创建一个包含用户 ID 的 Claims
	claims := jwt.BaseClaims{
		ID: id,
	}
	token, err := a.authService.CreateToken(claims)
	if err != nil {
		response.FailWithMessage("刷新token失败", c)
		return
	}
	response.OkWithData(token, c)
}

// register 处理用户注册请求
func (a *AuthController) register(c *gin.Context) {
	var userReq request.User
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

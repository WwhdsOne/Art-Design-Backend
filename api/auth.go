package api

import (
	"Art-Design-Backend/model/request"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/response"
	"Art-Design-Backend/pkg/utils"
	"Art-Design-Backend/service"
	"github.com/gin-gonic/gin"
)

func InitSecuredAuthRouter(r *gin.RouterGroup) {
	securedRouter := r.Group("/auth")
	securedRouter.POST("/logout", logout)
	securedRouter.POST("/refreshToken", refreshToken)
}

func InitOpenAuthRouter(r *gin.RouterGroup) {
	openRouter := r.Group("/auth")
	openRouter.POST("/login", login)
	openRouter.POST("/register", register)
}

// login 处理用户登录请求
// 参数: c *gin.Context - Gin框架的上下文对象，用于处理HTTP请求和响应
func login(c *gin.Context) {
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
	u, err := service.Login(c, loginReq)
	// 如果登录验证失败，返回错误响应并结束函数执行
	if err != nil {
		c.Error(err)
		return
	}
	token, err := jwt.CreateToken(jwt.BaseClaims{
		ID: u.ID,
	})
	if err != nil {
		response.FailWithMessage("生成token失败", c)
		return
	}
	// 返回生成的token
	response.OkWithData(token, c)
}

func logout(c *gin.Context) {
	// 从请求头中获取 token
	token := utils.GetToken(c)
	// 调用 jwt 包中的 LogoutToken 函数注销 token
	err := jwt.LogoutToken(token)
	if err != nil {
		response.FailWithMessage("注销失败", c)
		return
	}
	response.OkWithMessage("注销成功", c)
}

// refreshToken 处理用户刷新 token 请求
func refreshToken(c *gin.Context) {
	// 从上下文中获取用户 ID
	id := utils.GetUserID(c)
	token, err := jwt.CreateToken(jwt.BaseClaims{
		ID: id,
	})
	if err != nil {
		response.FailWithMessage("刷新token失败", c)
		return
	}
	response.OkWithData(token, c)
}

// register 处理用户注册请求
func register(c *gin.Context) {
	var userReq request.User
	err := c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		c.Set(gin.BindKey, userReq)
		return
	}
	err = service.AddUser(c, &userReq)
	if err != nil {
		c.Error(err)
		return
	}
	response.OkWithMessage("注册成功", c)
}

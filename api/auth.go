package api

import (
	"Art-Design-Backend/global"
	"Art-Design-Backend/model/request"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/redisx"
	"Art-Design-Backend/pkg/response"
	"Art-Design-Backend/pkg/utils"
	"Art-Design-Backend/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func InitSecuredAuthRouter(r *gin.RouterGroup) {
	r.POST("/logout", logout)
}

func InitOpenAuthRouter(r *gin.RouterGroup) {
	r.POST("/login", login)
	r.POST("/register", register)
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
	// 创建JWT claims，包含用户ID
	claim := global.JWT.CreateClaims(
		jwt.BaseClaims{
			ID: u.ID,
		},
	)
	// 将用户ID转换为字符串
	id := strconv.FormatInt(u.ID, 10)
	// 检查是否存在当前用户的会话
	session := redisx.Get(constant.SESSION + id)
	// 如果会话已存在，直接返回会话信息
	if session != "" {
		response.OkWithData(session, c)
		return
	}
	// 生成JWT token
	token, err := global.JWT.CreateToken(claim)
	// 如果生成token失败，返回错误响应并结束函数执行
	if err != nil {
		response.FailWithMessage("生成token失败", c)
		return
	}
	// 设置token方便获取是否登录
	err = redisx.Set(constant.LOGIN+token, id, global.JWT.ExpiresTime)
	// 如果设置token失败，返回错误响应并结束函数执行
	if err != nil {
		response.FailWithMessage("设置token失败", c)
		return
	}
	// 设置会话防止重复登录
	err = redisx.Set(constant.SESSION+id, token, global.JWT.ExpiresTime)
	// 如果设置会话失败，返回错误响应并结束函数执行
	if err != nil {
		response.FailWithMessage("设置会话失败", c)
		return
	}
	// 返回生成的token
	response.OkWithData(token, c)
}

func logout(c *gin.Context) {
	// 从请求头中获取 token
	token := utils.GetToken(c)
	if token == "" {
		response.FailWithMessage("当前未登录", c)
		return
	}
	// 从 Redis 中获取用户 ID
	claims, err := global.JWT.ParseToken(token)

	if err != nil {
		response.FailWithMessage("解析token失败", c)
		return
	}

	id := strconv.FormatInt(claims.BaseClaims.ID, 10)

	// 删除 Redis 中的会话信息
	err = redisx.Delete(constant.SESSION + id)
	if err != nil {
		response.FailWithMessage("删除会话失败", c)
		return
	}

	// 删除 Redis 中的登录信息
	err = redisx.Delete(constant.LOGIN + token)
	if err != nil {
		response.FailWithMessage("删除登录信息失败", c)
		return
	}

	response.OkWithMessage("注销成功", c)
}

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

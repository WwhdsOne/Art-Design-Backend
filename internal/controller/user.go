package controller

import (
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/result"
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService // 创建一个AuthService实例
}

func NewUserController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.UserService) *UserController {
	userCtrl := &UserController{
		userService: service,
	}
	r := engine.Group("/api").Group("/user")
	r.Use(middleware.AuthMiddleware())
	{
		// 私有路由组（需要 JWT 认证）
		r.GET("/info", userCtrl.getUserInfo)
		r.POST("/page", userCtrl.getUserPage)
		r.POST("/update", userCtrl.updateUserBaseInfo)
		r.POST("/updatePassword", userCtrl.updateUserPassword)
		r.POST("/uploadAvatar", userCtrl.uploadAvatar)
	}
	return userCtrl
}

func (u *UserController) getUserInfo(c *gin.Context) {
	user, err := u.userService.GetUserById(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(user, c)
}

func (u *UserController) getUserPage(c *gin.Context) {
	var userQuery query.User
	if err := c.ShouldBindJSON(&userQuery); err != nil {
		_ = c.Error(err)
		return
	}
	userPageData, err := u.userService.GetUserPage(c, &userQuery)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithData(userPageData, c)
}

func (u *UserController) updateUserBaseInfo(c *gin.Context) {
	var userReq request.User
	err := c.ShouldBindJSON(&userReq)
	if err != nil {
		_ = c.Error(err)
		c.Set(gin.BindKey, userReq)
		return
	}
	err = u.userService.UpdateUserBaseInfo(c, &userReq)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("更新用户成功", c)
}

func (u *UserController) updateUserPassword(c *gin.Context) {
	var changePwd request.ChangePassword
	err := c.ShouldBindJSON(&changePwd)
	if err != nil {
		_ = c.Error(err)
		c.Set(gin.BindKey, changePwd)
		return
	}
	err = u.userService.UpdateUserPassword(c, &changePwd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result.OkWithMessage("更新密码成功", c)
}

func (u *UserController) uploadAvatar(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		result.FailWithMessage("请选择要上传的文件", c)
		return
	}

	// 打开上传的文件流
	src, err := file.Open()
	if err != nil {
		result.FailWithMessage("无法打开上传的文件", c)
		return
	}
	defer src.Close()

	// 检查文件大小是否超过 2MB
	if file.Size > 2<<20 { // 2 MB
		result.FailWithMessage("文件大小不能超过 2MB", c)
		return
	}

	// 调用 service 层处理上传逻辑
	avatarURL, err := u.userService.UploadAvatar(c, file.Filename, src)
	if err != nil {
		result.FailWithMessage("头像上传失败: "+err.Error(), c)
		return
	}

	result.OkWithData(avatarURL, c)
}

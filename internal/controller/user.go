package controller

import (
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/service"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

//
//import (
//	"Art-Design-Backend/model/request"
//	"Art-Design-Backend/pkg/auth"
//	"Art-Design-Backend/pkg/response"
//	"Art-Design-Backend/pkg/utils"
//	"Art-Design-Backend/service"
//	"github.com/gin-gonic/gin"
//)
//
//func InitUserRouter(r *gin.RouterGroup) {
//	userRouter := r.Group("/user")
//	//userRouter.DELETE("/delete/:ids", deleteUser)
//	//userRouter.PUT("/update/:id", updateUser)
//	//userRouter.POST("/page", userPage)
//	userRouter.GET("/info", getUserInfo)
//}

type UserController struct {
	userService *service.UserService // 创建一个AuthService实例
}

func NewUserController(engine *gin.Engine, middleware *middleware.Middlewares, service *service.UserService) *UserController {
	userCtrl := &UserController{
		userService: service,
	}
	r := engine.Group("/api").Group("/user").Use(middleware.AuthMiddleware())
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
		c.Error(err)
		return
	}
	response.OkWithData(user, c)
}

func (u *UserController) getUserPage(c *gin.Context) {
	var userQuery query.User
	if err := c.ShouldBindJSON(&userQuery); err != nil {
		c.Error(err)
		return
	}
	user, err := u.userService.GetUserPage(c, &userQuery)
	if err != nil {
		c.Error(err)
		return
	}
	response.OkWithData(user, c)
}

//	func deleteUser(c *gin.Context) {
//		ids, err := utils.ParseIDs(c)
//		if err != nil {
//			response.FailWithMessage("参数错误", c)
//			return
//		}
//		err = service.DeleteUser(ids, auth.GetUserID(c))
//		if err != nil {
//			response.FailWithMessage("用户删除失败", c)
//			return
//		}
//		response.OkWithMessage("删除用户成功", c)
//	}
func (u *UserController) updateUserBaseInfo(c *gin.Context) {
	var userReq request.User
	err := c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		c.Set(gin.BindKey, userReq)
		return
	}
	err = u.userService.UpdateUserBaseInfo(c, &userReq)
	if err != nil {
		c.Error(err)
		return
	}
	response.OkWithMessage("更新用户成功", c)
}

func (u *UserController) updateUserPassword(c *gin.Context) {
	var changePwd request.ChangePassword
	err := c.ShouldBindJSON(&changePwd)
	if err != nil {
		c.Error(err)
		c.Set(gin.BindKey, changePwd)
		return
	}
	err = u.userService.UpdateUserPassword(c, &changePwd)
	if err != nil {
		c.Error(err)
		return
	}
	response.OkWithMessage("更新密码成功", c)
}

func (u *UserController) uploadAvatar(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		response.FailWithMessage("请选择要上传的文件", c)
		return
	}

	// 打开上传的文件流
	src, err := file.Open()
	if err != nil {
		response.FailWithMessage("无法打开上传的文件", c)
		return
	}
	defer src.Close()

	// 检查文件大小是否超过 2MB
	if file.Size > 2<<20 { // 2 MB
		response.FailWithMessage("文件大小不能超过 2MB", c)
		return
	}

	// 调用 service 层处理上传逻辑
	avatarURL, err := u.userService.UploadAvatar(c, file.Filename, src)
	if err != nil {
		response.FailWithMessage("头像上传失败: "+err.Error(), c)
		return
	}

	response.OkWithData(avatarURL, c)
}

//
////	func userPage(c *gin.Context) {
////		var user query.User
////		err := c.ShouldBindJSON(&user)
////		if err != nil {
////			response.FailWithMessage("分页参数填写错误", c)
////			return
////		}
////		users, total, err := service.UserPage(&user)
////		if err != nil {
////			response.FailWithMessage("分页查询失败", c)
////			return
////		}
////		pageResp := base.BuildPageResp[resp.User](users, total, user.PaginationReq)
////		response.OkWithData(pageResp, c)
////	}

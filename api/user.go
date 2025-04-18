package api

import (
	"Art-Design-Backend/model/request"
	"Art-Design-Backend/pkg/response"
	"Art-Design-Backend/pkg/utils"
	"Art-Design-Backend/service"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	userRouter := r.Group("/user")
	//userRouter.DELETE("/delete/:ids", deleteUser)
	//userRouter.PUT("/update/:id", updateUser)
	//userRouter.POST("/page", userPage)
	userRouter.GET("/info", getUserInfo)
}

func deleteUser(c *gin.Context) {
	ids, err := utils.ParseIDs(c)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = service.DeleteUser(ids, utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage("用户删除失败", c)
		return
	}
	response.OkWithMessage("删除用户成功", c)
}

func updateUser(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	var useReq request.User
	err = c.ShouldBindJSON(&useReq)
	if err != nil {
		response.FailWithMessage("用户名或密码不能为空", c)
		return
	}
	err = service.UpdateUser(c, &useReq, id)
	if err != nil {
		response.FailWithMessage("更新用户失败", c)
		return
	}
	response.OkWithMessage("更新用户成功", c)
}

//	func userPage(c *gin.Context) {
//		var user query.User
//		err := c.ShouldBindJSON(&user)
//		if err != nil {
//			response.FailWithMessage("分页参数填写错误", c)
//			return
//		}
//		users, total, err := service.UserPage(&user)
//		if err != nil {
//			response.FailWithMessage("分页查询失败", c)
//			return
//		}
//		pageResp := base.BuildPageResp[resp.User](users, total, user.PaginationReq)
//		response.OkWithData(pageResp, c)
//	}
func getUserInfo(c *gin.Context) {
	id := utils.GetUserID(c)
	userResp, err := service.GetUserById(id)
	if err != nil {
		c.Error(err)
		return
	}
	response.OkWithData(userResp, c)
}

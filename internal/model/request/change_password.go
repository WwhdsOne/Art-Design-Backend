package request

import "Art-Design-Backend/internal/model/common"

type ChangePassword struct {
	ID              common.LongStringID `json:"id" label:"用户ID" binding:"required"`
	OldPassword     string              `json:"oldPassword" label:"旧密码" binding:"required"`
	NewPassword     string              `json:"newPassword" label:"新密码" binding:"required,min=8,max=64,strongpassword"`
	ConfirmPassword string              `json:"confirmPassword" label:"确认密码" binding:"required,min=8,max=64,strongpassword,eqfield=NewPassword"`
}

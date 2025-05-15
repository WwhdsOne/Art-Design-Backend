package request

import (
	"Art-Design-Backend/internal/model/base"
)

type User struct {
	ID           base.LongStringID `json:"id" label:"ID" binding:"required"`
	RealName     string            `json:"realname" label:"真实姓名" binding:"max=50"`
	Nickname     string            `json:"nickname" label:"昵称" binding:"required,min=5,max=24"`
	Gender       int8              `json:"gender" label:"性别" binding:"oneof=1 2"`
	Email        string            `json:"email" label:"邮箱" binding:"omitempty,email,max=100"`
	Phone        string            `json:"phone" label:"手机号" binding:"omitempty,max=30"`
	Address      string            `json:"address" label:"地址" binding:"max=256"`
	Introduction string            `json:"introduction" label:"个人介绍" binding:"max=256"`
	Occupation   string            `json:"occupation" label:"职业" binding:"max=50"`
	Tags         []string          `json:"tags" label:"个人标签" binding:"dive,max=20"`
}

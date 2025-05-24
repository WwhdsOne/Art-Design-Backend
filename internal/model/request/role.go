package request

import "Art-Design-Backend/internal/model/base"

type Role struct {
	ID          base.LongStringID `json:"id" label:"角色ID"`
	Name        string            `json:"name" binding:"required,max=10" label:"角色名称"`
	Code        string            `json:"code" binding:"required,min=2,max=10" label:"角色编码"`
	Description string            `json:"description" label:"角色描述"`
	Status      int8              `json:"status" label:"角色状态"`
}

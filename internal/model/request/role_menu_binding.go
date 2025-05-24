package request

import "Art-Design-Backend/internal/model/base"

type RoleMenuBinding struct {
	RoleId  base.LongStringID  `json:"role_id" label:"角色ID" binding:"required"`
	MenuIds base.LongStringIDs `json:"menu_ids" label:"菜单ID列表" binding:"required"`
}

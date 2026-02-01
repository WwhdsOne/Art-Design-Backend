package request

import "Art-Design-Backend/internal/model/common"

type RoleMenuBinding struct {
	RoleID  common.LongStringID  `json:"role_id" label:"角色ID" binding:"required"`
	MenuIDs common.LongStringIDs `json:"menu_ids" label:"菜单ID列表" binding:"required"`
}

package request

import "Art-Design-Backend/internal/model/common"

type RoleMenuBinding struct {
	RoleId  common.LongStringID  `json:"role_id" label:"角色ID" binding:"required"`
	MenuIds common.LongStringIDs `json:"menu_ids" label:"菜单ID列表" binding:"required"`
}

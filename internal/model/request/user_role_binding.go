package request

import "Art-Design-Backend/internal/model/common"

type UserRoleBinding struct {
	UserID          common.LongStringID  `json:"user_id" label:"用户ID" binding:"required"`
	OriginalRoleIDs common.LongStringIDs `json:"original_role_ids" label:"原始角色ID列表" binding:"required"`
	RoleIDs         common.LongStringIDs `json:"role_ids" label:"角色ID列表" binding:"required"`
}

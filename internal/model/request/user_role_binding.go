package request

import "Art-Design-Backend/internal/model/base"

type UserRoleBinding struct {
	UserId          base.LongStringID  `json:"user_id" label:"用户ID" binding:"required"`
	OriginalRoleIds base.LongStringIDs `json:"original_role_ids" label:"原始角色ID列表" binding:"required"`
	RoleIds         base.LongStringIDs `json:"role_ids" label:"角色ID列表" binding:"required"`
}

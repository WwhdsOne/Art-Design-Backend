package response

import "Art-Design-Backend/internal/model/common"

type UserRoleBinding struct {
	Roles      []*SimpleRole        `json:"roles"`
	HasRoleIDs common.LongStringIDs `json:"has_role_ids"`
}

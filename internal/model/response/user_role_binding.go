package response

import "Art-Design-Backend/internal/model/base"

type UserRoleBinding struct {
	Roles      []*SimpleRole      `json:"roles"`
	HasRoleIDs base.LongStringIDs `json:"has_role_ids"`
}

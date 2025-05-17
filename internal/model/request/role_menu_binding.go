package request

import "Art-Design-Backend/internal/model/base"

type RoleMenuBinding struct {
	RoleId  base.LongStringID  `json:"role_id"`
	MenuIds base.LongStringIDs `json:"menu_ids"`
}

package response

import "Art-Design-Backend/internal/model/base"

type RoleMenuBinding struct {
	Menus      []*SimpleMenu      `json:"menus"`
	HasMenuIDs base.LongStringIDs `json:"has_menu_ids"`
}

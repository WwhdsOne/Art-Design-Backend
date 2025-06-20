package response

import "Art-Design-Backend/internal/model/common"

type RoleMenuBinding struct {
	Menus      []*SimpleMenu        `json:"menus"`
	HasMenuIDs common.LongStringIDs `json:"has_menu_ids"`
}

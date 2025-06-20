package request

import "Art-Design-Backend/internal/model/common"

type MenuAuth struct {
	ID       *common.LongStringID `json:"id" label:"菜单ID"`
	ParentID common.LongStringID  `json:"parentID" binding:"required" label:"父菜单ID"`
	Title    string               `json:"title" binding:"required" label:"菜单标题"`
	AuthCode string               `json:"authCode" binding:"required" label:"权限编码"`
	Sort     int                  `json:"sort" binding:"required" label:"排序"`
	Type     int8                 `json:"type" binding:"required" label:"菜单类型"`
}

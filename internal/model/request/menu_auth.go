package request

import "Art-Design-Backend/internal/model/base"

type MenuAuth struct {
	ParentID base.LongStringID `json:"parentID"`
	Title    string            `json:"title"`
	Code     string            `json:"code"`
	Sort     int               `json:"sort"`
	Type     int8              `json:"type"`
}

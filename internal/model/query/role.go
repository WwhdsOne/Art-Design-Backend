package query

import "Art-Design-Backend/internal/model/base"

type Role struct {
	Name string `json:"name"`
	base.PaginationReq
}

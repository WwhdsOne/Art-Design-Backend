package query

import "Art-Design-Backend/internal/model/common"

type Role struct {
	Name string `json:"name"`
	common.PaginationReq
}

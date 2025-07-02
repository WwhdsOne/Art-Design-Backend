package query

import "Art-Design-Backend/internal/model/common"

type AIAgent struct {
	common.PaginationReq
	Name *string `json:"name"`
}

package query

import "Art-Design-Backend/internal/model/common"

type AIProvider struct {
	common.PaginationReq

	Name *string `json:"name"`

	Enabled *bool `json:"enabled"`
}

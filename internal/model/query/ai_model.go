package query

import (
	"Art-Design-Backend/internal/model/common"
)

type AIModel struct {
	common.PaginationReq
	Model     *string `json:"model"`
	Provider  *string `json:"provider"`
	Enabled   *bool   `json:"enabled"`
	ModelType *string `json:"model_type"` // chat / embedding / multimodal
}

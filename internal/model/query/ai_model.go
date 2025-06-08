package query

import (
	"Art-Design-Backend/internal/model/base"
)

type AIModel struct {
	base.PaginationReq
	Model     *string `json:"model"`
	Provider  *string `json:"provider"`
	Enabled   *bool   `json:"enabled"`
	ModelType *string `json:"model_type"` // chat / embedding / multimodal
}

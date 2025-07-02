package request

import "Art-Design-Backend/internal/model/common"

type AIAgent struct {
	ID           common.LongStringID `json:"id"`
	Name         string              `json:"name" binding:"required"`
	Description  string              `json:"description"`
	ModelID      common.LongStringID `json:"modelID" binding:"required"`
	SystemPrompt string              `json:"systemPrompt" binding:"required"`
}

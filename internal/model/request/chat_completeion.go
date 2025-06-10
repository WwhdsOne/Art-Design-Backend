package request

import "Art-Design-Backend/internal/model/base"

type ChatCompletion struct {
	ID     base.LongStringID `json:"id" binding:"required" label:"模型ID"`
	Prompt string            `json:"prompt" binding:"required" label:"提示词"`
}

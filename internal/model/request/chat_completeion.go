package request

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/pkg/ai"
)

type ChatCompletion struct {
	ID       common.LongStringID `json:"id" binding:"required" label:"模型ID"`
	Messages []ai.ChatMessage    `json:"messages" binding:"required" label:"消息列表"`
}

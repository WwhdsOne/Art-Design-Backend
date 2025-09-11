package request

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/pkg/ai"
)

type ChatCompletion struct {
	ID              common.LongStringID `json:"id" binding:"required" label:"模型ID"`
	Messages        []ai.ChatMessage    `json:"messages" binding:"required" label:"消息列表"`
	ConversationID  common.LongStringID `json:"conversation_id" binding:"omitempty" label:"会话ID"`
	KnowledgeBaseID common.LongStringID `json:"knowledge_base_id" binding:"omitempty" label:"关联知识库ID"`
	Files           []string            `json:"files" binding:"omitempty" label:"上传文件"`
}

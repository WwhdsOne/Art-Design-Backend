package request

import "Art-Design-Backend/internal/model/common"

type CreateConversationRequest struct {
	Title string `json:"title" binding:"required,max=100"`
}

type RenameConversationRequest struct {
	ID    common.LongStringID `json:"id" binding:"required"`
	Title string              `json:"title" binding:"required,max=100"`
}

type DeleteConversationRequest struct {
	ID common.LongStringID `json:"id" binding:"required"`
}

type CreateMessageRequest struct {
	ConversationID int64  `json:"conversation_id,string" binding:"required"`
	Content        string `json:"content" binding:"required"`
}

type GetMessagesRequest struct {
	ConversationID int64 `form:"conversation_id" binding:"required"`
}

type CreateActionRequest struct {
	MessageID  int64   `json:"message_id,string" binding:"required"`
	ActionType string  `json:"action_type" binding:"required,oneof=goto click input select scroll wait"`
	Sequence   int     `json:"sequence"`
	URL        *string `json:"url,omitempty"`
	Selector   *string `json:"selector,omitempty"`
	Value      *string `json:"value,omitempty"`
	Distance   *int    `json:"distance,omitempty"`
	Timeout    *int    `json:"timeout,omitempty"`
}

type GetActionsRequest struct {
	MessageID int64 `form:"message_id" binding:"required"`
}

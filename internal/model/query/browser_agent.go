package query

import "Art-Design-Backend/internal/model/common"

type BrowserAgentConversation struct {
	Title string `json:"title" form:"title"`
	State string `json:"state" form:"state"`
	common.PaginationReq
}

type BrowserAgentMessage struct {
	ConversationID int64  `json:"conversation_id" form:"conversation_id"`
	State          string `json:"state" form:"state"`
	common.PaginationReq
}

type BrowserAgentAction struct {
	MessageID  int64  `json:"message_id" form:"message_id"`
	ActionType string `json:"action_type" form:"action_type"`
	Status     string `json:"status" form:"status"`
	common.PaginationReq
}

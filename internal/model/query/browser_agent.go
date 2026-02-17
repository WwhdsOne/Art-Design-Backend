package query

import "Art-Design-Backend/internal/model/common"

type BrowserAgentConversation struct {
	Title string `json:"title" form:"title"`
	State string `json:"state" form:"state"`
	common.PaginationReq
}

type BrowserAgentMessage struct {
	ConversationID common.LongStringID `json:"conversation_id" form:"conversation_id"`
	State          string              `json:"state" form:"state"`
	common.PaginationReq
}

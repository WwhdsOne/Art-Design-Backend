package request

import "Art-Design-Backend/internal/model/common"

type KnowledgeBase struct {
	ID      common.LongStringID  `json:"id"`
	Name    string               `json:"name"`
	FileIDs common.LongStringIDs `json:"file_ids"`
}

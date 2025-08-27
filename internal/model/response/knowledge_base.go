package response

import "github.com/dromara/carbon/v2"

type KnowledgeBase struct {
	ID          int64           `json:"id,string"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CreatedAt   carbon.DateTime `json:"created_at"`
	UpdatedAt   carbon.DateTime `json:"updated_at"`
}

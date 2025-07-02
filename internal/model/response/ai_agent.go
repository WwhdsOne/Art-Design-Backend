package response

type AIAgent struct {
	ID           int64  `json:"id,string"`
	Name         string `json:"name" `
	Description  string `json:"description"`
	ModelID      int64  `json:"modelID,string" `
	SystemPrompt string `json:"systemPrompt" `
}

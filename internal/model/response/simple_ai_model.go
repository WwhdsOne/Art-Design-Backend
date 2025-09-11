package response

type SimpleAIModel struct {
	ID        int64  `json:"id,string"`
	Model     string `json:"model"`
	Icon      string `json:"icon"`
	ModelType string `json:"model_type"`
}

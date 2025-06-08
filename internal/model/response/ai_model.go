package response

import (
	"github.com/shopspring/decimal"
)

type AIModel struct {
	ID       int64  `json:"id,string"`
	Model    string `json:"model"`
	Provider string `json:"provider"`
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	ModelID  string `json:"model_id"`
	Icon     string `json:"icon"`

	PricePromptPer1M       decimal.Decimal `json:"price_prompt_per_1m"`
	PricePromptCachedPer1M decimal.Decimal `json:"price_prompt_cached_per_1m"`
	PriceCompletionPer1M   decimal.Decimal `json:"price_completion_per_1m"`

	Currency         string `json:"currency"`
	Enabled          bool   `json:"enabled"`
	MaxContextTokens int    `json:"max_context_tokens"`
	ModelType        string `json:"model_type"`
}

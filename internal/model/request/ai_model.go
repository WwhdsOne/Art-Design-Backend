package request

import (
	"github.com/shopspring/decimal"
)

type AIModel struct {
	Model    string `json:"model" binding:"required" label:"模型名称"`
	Provider string `json:"provider" binding:"required" label:"模型提供商"`
	BaseURL  string `json:"base_url" binding:"required" label:"API 基础地址"`
	APIKey   string `json:"api_key" binding:"required" label:"API 密钥"`
	ModelID  string `json:"model_id" binding:"required" label:"模型接口标识"`
	Icon     string `json:"icon" label:"模型图标 URL"`

	PricePromptPer1M       decimal.Decimal `json:"price_prompt_per_1m" binding:"required" label:"输入 token 单价"`
	PricePromptCachedPer1M decimal.Decimal `json:"price_prompt_cached_per_1m" binding:"required" label:"命中缓存输入单价"`
	PriceCompletionPer1M   decimal.Decimal `json:"price_completion_per_1m" binding:"required" label:"输出 token 单价"`

	Currency         string `json:"currency" binding:"required" label:"币种"` // 如 CNY
	Enabled          bool   `json:"enabled" label:"是否启用"`
	MaxContextTokens int    `json:"max_context_tokens" binding:"required" label:"最大上下文长度"`
	ModelType        string `json:"model_type" binding:"required" label:"模型类型"` // chat / embedding / multimodal
}

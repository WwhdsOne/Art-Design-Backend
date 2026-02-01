package request

import (
	"Art-Design-Backend/internal/model/common"

	"github.com/shopspring/decimal"
)

type AIModel struct {
	ID         common.LongStringID `json:"id" label:"模型ID"`
	Model      string              `json:"model" binding:"required" label:"模型名称"`
	ProviderID common.LongStringID `json:"provider" binding:"required" label:"模型提供商"`
	APIPath    string              `json:"api_path" binding:"required" label:"API 路径"`
	ModelID    string              `json:"model_id" label:"模型官方ID"`
	Icon       string              `json:"icon" label:"模型图标 URL"`

	PricePromptPer1M       decimal.Decimal `json:"price_prompt_per_1m" binding:"required" label:"百万token输入价格"`
	PricePromptCachedPer1M decimal.Decimal `json:"price_prompt_cached_per_1m" label:"百万token输入(缓存)价格"`
	PriceCompletionPer1M   decimal.Decimal `json:"price_completion_per_1m" binding:"required" label:"百万token输出价格"`

	Currency          string `json:"currency" binding:"required" label:"币种"` // 如 CNY
	Enabled           bool   `json:"enabled" binding:"required" label:"是否启用"`
	MaxContextTokens  int    `json:"max_context_tokens" binding:"required" label:"最大上下文长度"`
	MaxGenerateTokens int    `json:"max_generate_tokens" binding:"required" label:"最大生成长度"`
	ModelType         string `json:"model_type" binding:"required" label:"模型类型"` // chat / embedding / multimodal
}

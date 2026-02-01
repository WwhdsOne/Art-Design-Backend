package response

import (
	"github.com/shopspring/decimal"
)

type AIModel struct {
	ID string `json:"id"` // 使用 string 是因为 json:",string" 表示序列化为字符串

	Model    string `json:"model"`    // 模型名称
	Provider string `json:"provider"` // 供应商名称（如果是 ID 则改为 provider_id）

	APIPath string `json:"api_path"` // API 路径（BaseURL 后的部分）
	ModelID string `json:"model_id"` // 模型接口标识（如 gpt-4-1106-preview）
	Icon    string `json:"icon"`     // 模型图标 URL

	PricePromptPer1M       decimal.Decimal `json:"price_prompt_per_1m"`        // 正常输入 token 单价
	PricePromptCachedPer1M decimal.Decimal `json:"price_prompt_cached_per_1m"` // 命中缓存输入单价
	PriceCompletionPer1M   decimal.Decimal `json:"price_completion_per_1m"`    // 输出 token 单价

	Currency          string `json:"currency"`            // 币种
	Enabled           bool   `json:"enabled"`             // 是否启用
	MaxContextTokens  int    `json:"max_context_tokens"`  // 最大上下文长度
	MaxGenerateTokens int    `json:"max_generate_tokens"` // 最大生成长度
	ModelType         string `json:"model_type"`          // 模型类型：chat / embedding / multimodal
}

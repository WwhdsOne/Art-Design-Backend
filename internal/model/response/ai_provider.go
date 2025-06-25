package response

import "github.com/dromara/carbon/v2"

type AIProvider struct {
	// 供应商ID
	ID int64 `json:"id,string"`

	// 供应商名称
	Name string `json:"name"`

	// 调用 API 的基础地址，如 "https://api.openai.com/v1"
	BaseURL string `json:"base_url"`

	// API 密钥，建议加密存储
	APIKey string `json:"api_key"`

	// 是否启用该供应商，false 时该供应商下所有模型不可用
	Enabled bool `json:"enabled"`

	// 最大请求速率限制(次/分钟)，0 表示没有并发限制
	MaxRateLimit int `json:"max_rate_limit"`

	CreatedAt carbon.DateTime `json:"created_at"`

	CreatedBy int64 `json:"created_by,string"`

	UpdatedAt carbon.DateTime `json:"updated_at"`

	UpdatedBy int64 `json:"updated_by,string"`
}

package request

import "Art-Design-Backend/internal/model/common"

type AIProvider struct {
	// 供应商ID
	ID common.LongStringID `json:"id" label:"供应商ID"`

	// 供应商名称，例如 openai、anthropic、cohere 等
	Name string `json:"name" label:"供应商名称" binding:"required,min=2,max=100"`

	// 调用 API 的基础地址，如 "https://api.openai.com/v1"
	BaseURL string `json:"base_url" label:"API基础地址" binding:"omitempty,url,max=200"`

	// API 密钥，建议加密存储
	APIKey string `json:"api_key" label:"API密钥"`

	// 是否启用该供应商，false 时该供应商下所有模型不可用
	Enabled bool `json:"enabled" label:"是否启用"`

	// 最大请求速率限制(次/分钟)，0 表示没有并发限制
	MaxRateLimit int `json:"max_rate_limit" label:"最大请求速率(次/分钟)" binding:"min=0"`
}

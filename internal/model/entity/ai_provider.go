package entity

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/pkg/constant/tablename"
)

// AIProvider 供应商表
type AIProvider struct {
	common.BaseModel

	// 供应商名称，例如 openai、anthropic、cohere 等
	Name string `gorm:"type:varchar(100);not null;unique;comment:供应商名称，例如 openai"`

	// 调用 API 的基础地址，如 "https://api.openai.com/v1"
	BaseURL string `gorm:"type:varchar(200);comment:调用 API 的基础地址"`

	// API 密钥，建议加密存储
	APIKey string `gorm:"type:varchar(100);comment:API 密钥，建议加密存储"`

	// 是否启用该供应商，false 时该供应商下所有模型不可用
	Enabled bool `gorm:"default:true;comment:是否启用该供应商"`

	// 最大请求速率限制(次/分钟)，0 表示没有并发限制
	MaxRateLimit int `gorm:"comment:最大请求速率限制(次/分钟)，0表示没有限制"`
}

func (a *AIProvider) TableName() string {
	return tablename.AIProviderTableName
}

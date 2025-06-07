package entity

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/pkg/constant/tablename"
	"github.com/shopspring/decimal"
)

type AIModel struct {
	base.BaseModel
	Model    string `gorm:"type:varchar(100);not null"` // 模型名称，如 "gpt-4"
	Provider string `gorm:"type:varchar(100);not null"` // 模型提供商，如 "openai"
	BaseURL  string `gorm:"type:varchar(200)"`          // 调用API基础地址
	APIKey   string `gorm:"type:varchar(100)"`          // API密钥，建议加密存储
	ModelID  string `gorm:"type:varchar(100)"`          // 模型接口标识，如 "gpt-4-1106-preview"
	Icon     string `gorm:"type:varchar(255)"`          // 模型图标URL

	PricePromptPer1M              decimal.Decimal `gorm:"type:numeric(15,8)"` // 正常输入 token 单价（每百万 token）
	PricePromptCachedPer1M        decimal.Decimal `gorm:"type:numeric(15,8)"` // 命中缓存输入 token 单价
	PriceCacheStoragePer1MPerHour decimal.Decimal `gorm:"type:numeric(15,8)"` // 缓存存储单价（每百万 token 每小时）
	PriceCompletionPer1M          decimal.Decimal `gorm:"type:numeric(15,8)"` // 输出 token 单价

	Currency         string `gorm:"type:varchar(10);default:'CNY'"` // 币种，默认人民币 CNY
	Enabled          bool   `gorm:"default:true"`                   // 是否启用
	MaxContextTokens int    `gorm:"not null"`                       // 最大上下文长度（单位：token）
}

func (a *AIModel) TableName() string {
	return tablename.AIModelTableName
}

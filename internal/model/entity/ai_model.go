package entity

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/pkg/constant/tablename"
	"github.com/shopspring/decimal"
)

type AIModel struct {
	base.BaseModel

	Model    string `gorm:"type:varchar(100);not null;comment:模型名称，例如 gpt-4"`
	Provider string `gorm:"type:varchar(100);not null;comment:模型提供商，例如 openai"`
	BaseURL  string `gorm:"type:varchar(200);comment:调用 API 的基础地址"`
	APIKey   string `gorm:"type:varchar(100);comment:API 密钥，建议加密存储"`
	ModelID  string `gorm:"type:varchar(100);comment:模型接口标识，例如 gpt-4-1106-preview"`
	Icon     string `gorm:"type:varchar(255);comment:模型图标的 URL 地址"`

	PricePromptPer1M       decimal.Decimal `gorm:"type:numeric(15,8);comment:正常输入 token 单价（每百万 token）"`
	PricePromptCachedPer1M decimal.Decimal `gorm:"type:numeric(15,8);comment:命中缓存输入 token 单价（每百万 token）"`
	PriceCompletionPer1M   decimal.Decimal `gorm:"type:numeric(15,8);comment:输出 token 单价（每百万 token）"`

	Currency         string `gorm:"type:varchar(10);default:'CNY';comment:计价币种，默认人民币 CNY"`
	Enabled          bool   `gorm:"default:true;comment:是否启用该模型"`
	MaxContextTokens int    `gorm:"not null;comment:最大上下文长度（单位：token）"`
	ModelType        string `gorm:"type:varchar(50);not null;comment:模型类型，如 chat、embedding、multimodal"`
}

func (a *AIModel) TableName() string {
	return tablename.AIModelTableName
}

package entity

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/pkg/constant/tablename"
)

type AIAgent struct {
	common.BaseModel
	Name         string `gorm:"type:varchar(100);not null;column:name"`
	Description  string `gorm:"type:text;column:description"`
	ModelID      int64  `gorm:"type:int8;not null;column:model_id"`
	SystemPrompt string `gorm:"type:text;column:system_prompt"`
}

func (a *AIAgent) TableName() string {
	return tablename.AIAgentTableName
}

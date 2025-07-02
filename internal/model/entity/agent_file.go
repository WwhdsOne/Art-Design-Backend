package entity

import (
	"Art-Design-Backend/internal/model/common"
)

type AgentFile struct {
	common.BaseModel
	AgentID int64  `gorm:"not null;index"`
	FileURL string `gorm:"size:500;not null"` // 文件路径（本地或对象存储）
}

func (a *AgentFile) TableName() string {
	return "agent_file"
}

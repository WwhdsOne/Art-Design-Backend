package entity

import (
	"Art-Design-Backend/pkg/constant/tablename"
	"time"
)

type BrowserAgentMessage struct {
	ID             int64     `gorm:"type:bigint;primaryKey;comment:雪花ID"`
	ConversationID int64     `gorm:"column:conversation_id;not null;index;comment:会话ID"`
	Content        string    `gorm:"column:content;type:text;comment:用户任务描述"`
	CreatedAt      time.Time `gorm:"type:timestamp;column:created_at;autoCreateTime"`
}

func (b *BrowserAgentMessage) TableName() string {
	return tablename.BrowserAgentMessageTableName
}

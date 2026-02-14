package entity

import (
	"Art-Design-Backend/pkg/constant/tablename"
)

// 消息角色常量
const (
	MessageRoleUser      = "user"
	MessageRoleAssistant = "assistant"
)

// BrowserAgentMessage 浏览器代理消息实体
type BrowserAgentMessage struct {
	ID             int64  `gorm:"type:bigint;primaryKey;comment:雪花ID"`
	ConversationID int64  `gorm:"column:conversation_id;not null;index;comment:会话ID"`
	Role           string `gorm:"column:role;type:varchar(20);not null;comment:角色(user/assistant)"`
	Content        string `gorm:"column:content;type:text;comment:消息内容"`
}

// TableName 指定消息表名
func (b *BrowserAgentMessage) TableName() string {
	return tablename.BrowserAgentMessageTableName
}

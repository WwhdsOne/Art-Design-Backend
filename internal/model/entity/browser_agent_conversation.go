package entity

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/pkg/constant/tablename"
)

// 会话状态常量
const (
	ConversationStateRunning             = "running"
	ConversationStateWaitingVerification = "waiting_verification"
	ConversationStateFinished            = "finished"
	ConversationStateError               = "error"
)

// BrowserAgentConversation 浏览器代理会话实体
type BrowserAgentConversation struct {
	common.BaseModel
	Title string `gorm:"column:title;type:varchar(100);not null;comment:会话标题"`
	State string `gorm:"column:state;type:varchar(30);default:running;comment:状态"`
}

// TableName 指定会话表名
func (b *BrowserAgentConversation) TableName() string {
	return tablename.BrowserAgentConversationTableName
}

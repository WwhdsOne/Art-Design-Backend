package entity

import (
	"Art-Design-Backend/pkg/constant/tablename"
	"time"
)

// 会话状态常量
const (
	MessageStateRunning  = "running"
	MessageStateFinished = "finished"
	MessageStateError    = "error"
)

type BrowserAgentMessage struct {
	ID                 int64      `gorm:"type:bigint;primaryKey;comment:雪花ID"`
	ConversationID     int64      `gorm:"column:conversation_id;not null;index;comment:会话ID"`
	Content            string     `gorm:"column:content;type:text;comment:用户任务描述"`
	State              string     `gorm:"column:state;type:varchar(30);default:running;comment:状态"`
	PageURL            *string    `gorm:"column:page_url;type:varchar(500);comment:任务开始时的页面URL"`
	ElementCount       *int       `gorm:"column:element_count;comment:页面可交互元素数量"`
	TotalSteps         *int       `gorm:"column:total_steps;comment:总执行步数"`
	TotalExecutionTime *int       `gorm:"column:total_execution_time;comment:总耗时(ms)"`
	LLMModel           *string    `gorm:"column:llm_model;type:varchar(100);comment:使用的LLM模型"`
	LLMTokenUsage      *string    `gorm:"column:llm_token_usage;type:jsonb;comment:累计token消耗统计"`
	FinishedAt         *time.Time `gorm:"column:finished_at;comment:任务完成时间"`
	CreatedAt          time.Time  `gorm:"type:timestamp;column:created_at;autoCreateTime"`
}

func (b *BrowserAgentMessage) TableName() string {
	return tablename.BrowserAgentMessageTableName
}

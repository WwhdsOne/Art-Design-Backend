package entity

import (
	"Art-Design-Backend/pkg/constant/tablename"
	"time"
)

// 操作状态常量
const (
	ActionStatusPending = "pending"
	ActionStatusRunning = "running"
	ActionStatusSuccess = "success"
	ActionStatusFailed  = "failed"
	ActionStatusSkipped = "skipped"
)

// BrowserAgentAction 浏览器代理操作实体
type BrowserAgentAction struct {
	ID            int64     `gorm:"type:bigint;primaryKey;comment:雪花ID"`
	MessageID     int64     `gorm:"column:message_id;not null;index;comment:消息ID"`
	ActionType    string    `gorm:"column:action_type;type:varchar(30);not null;comment:操作类型(goto/click/input/select/scroll/wait)"`
	ElementIndex  *int      `gorm:"column:element_index;comment:元素编号索引"`
	Status        string    `gorm:"column:status;type:varchar(20);default:pending;comment:状态"`
	URL           *string   `gorm:"column:url;type:varchar(500);comment:URL(goto)"`
	Selector      *string   `gorm:"column:selector;type:varchar(500);comment:选择器(click/input/select)"`
	Value         *string   `gorm:"column:value;type:text;comment:值(input/select)"`
	Distance      *int      `gorm:"column:distance;comment:滚动距离"`
	Timeout       *int      `gorm:"column:timeout;comment:等待时间"`
	ErrorMessage  *string   `gorm:"column:error_message;type:text;comment:错误信息"`
	ExecutionTime *int      `gorm:"column:execution_time;comment:执行耗时(毫秒)"`
	PageURL       *string   `gorm:"column:page_url;type:varchar(500);comment:执行动作时的页面URL"`
	ElementText   *string   `gorm:"column:element_text;type:varchar(200);comment:操作元素的可见文本"`
	LLMThinking   *string   `gorm:"column:llm_thinking;type:text;comment:LLM思维链"`
	LLMResponse   *string   `gorm:"column:llm_response;type:text;comment:LLM原始返回JSON"`
	LLMTokenUsage *string   `gorm:"column:llm_token_usage;type:jsonb;comment:单步token消耗"`
	IsVision      *bool     `gorm:"column:is_vision;comment:是否经过视觉模型回退"`
	CreatedAt     time.Time `gorm:"type:timestamp;column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"type:timestamp;column:updated_at;autoUpdateTime"`
}

// TableName 指定操作表名
func (b *BrowserAgentAction) TableName() string {
	return tablename.BrowserAgentActionTableName
}

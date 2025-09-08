package entity

import (
	"Art-Design-Backend/pkg/constant/tablename"
	"time"

	"github.com/lib/pq"
)

type Message struct {
	ID              int64         `gorm:"type:bigint;column:id;primaryKey;autoIncrement:false"`
	ConversationID  int64         `gorm:"not null;index;comment:会话ID"`
	Role            string        `gorm:"type:varchar(20);not null;check:role IN ('user','assistant');comment:消息角色"`
	Content         string        `gorm:"type:text;not null;comment:消息内容"`
	FileChunkIDs    pq.Int64Array `gorm:"type:bigint[];comment:关联知识片段ID数组" `
	KnowledgeBaseID *int64        `gorm:"comment:知识库ID(可为空)"`
	CreatedAt       time.Time     `gorm:"not null;autoCreateTime;comment:创建时间"`
}

func (m *Message) TableName() string {
	return tablename.MessageTableName
}

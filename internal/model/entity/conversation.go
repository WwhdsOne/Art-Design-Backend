package entity

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/pkg/constant/tablename"
)

type Conversation struct {
	common.BaseModel
	Title string `gorm:"type:varchar(50);comment:标题"`
}

func (c *Conversation) TableName() string {
	return tablename.ConversationTableName
}

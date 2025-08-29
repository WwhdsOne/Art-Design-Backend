package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"

	"gorm.io/gorm"
)

type MessageDB struct {
	db *gorm.DB
}

func NewMessageDB(db *gorm.DB) *MessageDB {
	return &MessageDB{
		db: db,
	}
}

func (m *MessageDB) CreateMessage(ctx context.Context, e *entity.Message) error {
	if err := DB(ctx, m.db).Create(e).Error; err != nil {
		return errors.WrapDBError(err, "创建消息失败")
	}
	return nil
}
func (m *MessageDB) GetMessageByConversationID(ctx context.Context, id int64) (messages []*entity.Message, err error) {
	if err := DB(ctx, m.db).Where("conversation_id = ?", id).Find(&messages).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询会话消息失败")
	}
	return messages, nil
}

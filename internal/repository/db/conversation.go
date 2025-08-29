package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"

	"gorm.io/gorm"
)

type ConversationDB struct {
	db *gorm.DB
}

func NewConversationDB(db *gorm.DB) *ConversationDB {
	return &ConversationDB{
		db: db,
	}
}

func (c *ConversationDB) CreateConversation(ctx context.Context, e *entity.Conversation) error {
	if err := DB(ctx, c.db).Create(e).Error; err != nil {
		return errors.WrapDBError(err, "创建会话失败")
	}
	return nil
}

func (c *ConversationDB) GetHistoryConversation(ctx context.Context, userID int64) (res []*entity.Conversation, err error) {
	if err = DB(ctx, c.db).Select("id", "title", "created_at").
		Order("created_at DESC").
		Where("created_by = ?", userID).Find(&res).Error; err != nil {
		err = errors.WrapDBError(err, "获取历史会话失败")
		return
	}
	return
}

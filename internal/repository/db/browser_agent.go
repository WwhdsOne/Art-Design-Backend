package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	myerrors "Art-Design-Backend/pkg/errors"
	"context"
	"errors"

	"gorm.io/gorm"
)

type BrowserAgentDB struct {
	db *gorm.DB
}

func NewBrowserAgentDB(db *gorm.DB) *BrowserAgentDB {
	return &BrowserAgentDB{db: db}
}

// =========================
// Conversation CRUD
// =========================

func (r *BrowserAgentDB) CreateConversation(ctx context.Context, conv *entity.BrowserAgentConversation) error {
	if err := r.db.WithContext(ctx).Create(conv).Error; err != nil {
		return myerrors.WrapDBError(err, "创建浏览器智能体会话失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetConversationByID(ctx context.Context, id int64) (*entity.BrowserAgentConversation, error) {
	var conv entity.BrowserAgentConversation
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&conv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDBError("会话不存在")
		}
		return nil, myerrors.WrapDBError(err, "查询浏览器智能体会话失败")
	}
	return &conv, nil
}

func (r *BrowserAgentDB) ListConversationsByUserID(ctx context.Context, userID int64, queryParam *query.BrowserAgentConversation) ([]*entity.BrowserAgentConversation, int64, error) {
	var conversations []*entity.BrowserAgentConversation
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.BrowserAgentConversation{}).Where("created_by = ?", userID)

	if queryParam.Title != "" {
		db = db.Where("title LIKE ?", "%"+queryParam.Title+"%")
	}
	if queryParam.State != "" {
		db = db.Where("state = ?", queryParam.State)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, myerrors.WrapDBError(err, "统计浏览器智能体会话数量失败")
	}

	if err := db.Scopes(queryParam.Paginate()).Order("created_at DESC").Find(&conversations).Error; err != nil {
		return nil, 0, myerrors.WrapDBError(err, "查询浏览器智能体会话列表失败")
	}

	return conversations, total, nil
}

func (r *BrowserAgentDB) UpdateConversation(ctx context.Context, conv *entity.BrowserAgentConversation) error {
	if err := r.db.WithContext(ctx).Save(conv).Error; err != nil {
		return myerrors.WrapDBError(err, "更新浏览器智能体会话失败")
	}
	return nil
}

func (r *BrowserAgentDB) DeleteConversation(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&entity.BrowserAgentConversation{}, id).Error; err != nil {
		return myerrors.WrapDBError(err, "删除浏览器智能体会话失败")
	}
	return nil
}

func (r *BrowserAgentDB) UpdateConversationState(ctx context.Context, id int64, state string) error {
	if err := r.db.WithContext(ctx).Model(&entity.BrowserAgentConversation{}).
		Where("id = ?", id).Update("state", state).Error; err != nil {
		return myerrors.WrapDBError(err, "更新会话状态失败")
	}
	return nil
}

// =========================
// Message CRUD
// =========================

func (r *BrowserAgentDB) CreateMessage(ctx context.Context, msg *entity.BrowserAgentMessage) error {
	if err := r.db.WithContext(ctx).Create(msg).Error; err != nil {
		return myerrors.WrapDBError(err, "创建浏览器智能体消息失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetMessageByID(ctx context.Context, id int64) (*entity.BrowserAgentMessage, error) {
	var msg entity.BrowserAgentMessage
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&msg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDBError("消息不存在")
		}
		return nil, myerrors.WrapDBError(err, "查询浏览器智能体消息失败")
	}
	return &msg, nil
}

func (r *BrowserAgentDB) ListMessagesByConversationID(ctx context.Context, conversationID int64) ([]*entity.BrowserAgentMessage, error) {
	var messages []*entity.BrowserAgentMessage

	// 只按 conversation_id 查询
	db := r.db.WithContext(ctx).Model(&entity.BrowserAgentMessage{}).Where("conversation_id = ?", conversationID)

	// 查询所有消息
	if err := db.Find(&messages).Error; err != nil {
		return nil, myerrors.WrapDBError(err, "查询浏览器智能体消息列表失败")
	}

	return messages, nil
}

func (r *BrowserAgentDB) GetRecentMessages(ctx context.Context, conversationID int64, limit int) ([]*entity.BrowserAgentMessage, error) {
	var messages []*entity.BrowserAgentMessage
	err := r.db.WithContext(ctx).Model(&entity.BrowserAgentMessage{}).
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, myerrors.WrapDBError(err, "查询最近的浏览器智能体消息失败")
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, nil
}

// =========================
// Action CRUD
// =========================

func (r *BrowserAgentDB) CreateAction(ctx context.Context, action *entity.BrowserAgentAction) error {
	if err := r.db.WithContext(ctx).Create(action).Error; err != nil {
		return myerrors.WrapDBError(err, "创建浏览器智能体操作失败")
	}
	return nil
}

func (r *BrowserAgentDB) CreateActions(ctx context.Context, actions []*entity.BrowserAgentAction) error {
	if len(actions) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(&actions).Error; err != nil {
		return myerrors.WrapDBError(err, "批量创建浏览器智能体操作失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetActionByID(ctx context.Context, id int64) (*entity.BrowserAgentAction, error) {
	var action entity.BrowserAgentAction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&action).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDBError("操作不存在")
		}
		return nil, myerrors.WrapDBError(err, "查询浏览器智能体操作失败")
	}
	return &action, nil
}

func (r *BrowserAgentDB) ListActionsByMessageID(ctx context.Context, messageID int64) ([]*entity.BrowserAgentAction, error) {
	var actions []*entity.BrowserAgentAction

	// 只按 message_id 查询
	if err := r.db.WithContext(ctx).
		Model(&entity.BrowserAgentAction{}).
		Where("message_id = ?", messageID).
		Find(&actions).Error; err != nil {
		return nil, myerrors.WrapDBError(err, "查询浏览器智能体操作列表失败")
	}

	return actions, nil
}

func (r *BrowserAgentDB) UpdateAction(ctx context.Context, action *entity.BrowserAgentAction) error {
	if err := r.db.WithContext(ctx).Save(action).Error; err != nil {
		return myerrors.WrapDBError(err, "更新浏览器智能体操作失败")
	}
	return nil
}

func (r *BrowserAgentDB) UpdateActionStatus(ctx context.Context, id int64, status string, errMsg *string, execTime *int) error {
	updates := map[string]interface{}{"status": status}
	if errMsg != nil {
		updates["error_message"] = *errMsg
	}
	if execTime != nil {
		updates["execution_time"] = *execTime
	}

	if err := r.db.WithContext(ctx).Model(&entity.BrowserAgentAction{}).
		Where("id = ?", id).Updates(updates).Error; err != nil {
		return myerrors.WrapDBError(err, "更新操作状态失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetPendingActionsByMessageID(ctx context.Context, messageID int64) ([]*entity.BrowserAgentAction, error) {
	var actions []*entity.BrowserAgentAction
	err := r.db.WithContext(ctx).Model(&entity.BrowserAgentAction{}).
		Where("message_id = ? AND status = ?", messageID, entity.ActionStatusPending).
		Order("sequence ASC").
		Find(&actions).Error
	if err != nil {
		return nil, myerrors.WrapDBError(err, "查询待执行操作失败")
	}
	return actions, nil
}

// =========================
// 聚合查询
// =========================

func (r *BrowserAgentDB) GetMessageWithActions(ctx context.Context, messageID int64) (*entity.BrowserAgentMessage, []*entity.BrowserAgentAction, error) {
	msg, err := r.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, nil, err
	}

	var actions []*entity.BrowserAgentAction
	err = r.db.WithContext(ctx).Model(&entity.BrowserAgentAction{}).
		Where("message_id = ?", messageID).
		Order("sequence ASC").
		Find(&actions).Error
	if err != nil {
		return nil, nil, myerrors.WrapDBError(err, "查询消息操作列表失败")
	}

	return msg, actions, nil
}

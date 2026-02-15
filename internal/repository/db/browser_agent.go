package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/errors"
	"context"

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
	if err := DB(ctx, r.db).Create(conv).Error; err != nil {
		return errors.WrapDBError(err, "创建浏览器智能体会话失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetConversationByID(ctx context.Context, id int64) (conv *entity.BrowserAgentConversation, err error) {
	if err = DB(ctx, r.db).Where("id = ?", id).First(&conv).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询浏览器智能体会话失败")
	}
	return
}

func (r *BrowserAgentDB) ListConversationsByUserID(ctx context.Context, userID int64, queryParam *query.BrowserAgentConversation) (conversations []*entity.BrowserAgentConversation, total int64, err error) {
	db := DB(ctx, r.db).Model(&entity.BrowserAgentConversation{}).Where("created_by = ?", userID)

	if queryParam.Title != "" {
		db = db.Where("title LIKE ?", "%"+queryParam.Title+"%")
	}
	if queryParam.State != "" {
		db = db.Where("state = ?", queryParam.State)
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapDBError(err, "统计浏览器智能体会话数量失败")
	}

	if err = db.Scopes(queryParam.Paginate()).Order("created_at DESC").Find(&conversations).Error; err != nil {
		return nil, 0, errors.WrapDBError(err, "查询浏览器智能体会话列表失败")
	}

	return
}

func (r *BrowserAgentDB) UpdateConversation(ctx context.Context, conv *entity.BrowserAgentConversation) error {
	if err := DB(ctx, r.db).Save(conv).Error; err != nil {
		return errors.WrapDBError(err, "更新浏览器智能体会话失败")
	}
	return nil
}

func (r *BrowserAgentDB) DeleteConversation(ctx context.Context, id int64) error {
	if err := DB(ctx, r.db).Delete(&entity.BrowserAgentConversation{}, id).Error; err != nil {
		return errors.WrapDBError(err, "删除浏览器智能体会话失败")
	}
	return nil
}

func (r *BrowserAgentDB) DeleteMessagesByConversationID(ctx context.Context, conversationID int64) error {
	if err := DB(ctx, r.db).Where("conversation_id = ?", conversationID).Delete(&entity.BrowserAgentMessage{}).Error; err != nil {
		return errors.WrapDBError(err, "删除浏览器智能体消息失败")
	}
	return nil
}

func (r *BrowserAgentDB) DeleteActionsByMessageIDList(ctx context.Context, messageIDList []int64) error {
	if err := DB(ctx, r.db).
		Where("message_id IN (?)", messageIDList). // 关键修改：= → IN
		Delete(&entity.BrowserAgentAction{}).Error; err != nil {
		return errors.WrapDBError(err, "批量删除浏览器智能体操作失败")
	}
	return nil
}

func (r *BrowserAgentDB) UpdateConversationState(ctx context.Context, id int64, state string) error {
	if err := DB(ctx, r.db).Model(&entity.BrowserAgentConversation{}).
		Where("id = ?", id).Update("state", state).Error; err != nil {
		return errors.WrapDBError(err, "更新会话状态失败")
	}
	return nil
}

// =========================
// Message CRUD
// =========================

func (r *BrowserAgentDB) CreateMessage(ctx context.Context, msg *entity.BrowserAgentMessage) error {
	if err := DB(ctx, r.db).Create(msg).Error; err != nil {
		return errors.WrapDBError(err, "创建浏览器智能体消息失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetMessageByID(ctx context.Context, id int64) (msg *entity.BrowserAgentMessage, err error) {
	if err = DB(ctx, r.db).Where("id = ?", id).First(&msg).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errors.NewDBError("消息不存在")
		}
		return nil, errors.WrapDBError(err, "查询浏览器智能体消息失败")
	}
	return
}

func (r *BrowserAgentDB) ListMessagesByConversationID(ctx context.Context, conversationID int64) (messages []*entity.BrowserAgentMessage, err error) {
	if err = DB(ctx, r.db).Model(&entity.BrowserAgentMessage{}).
		Where("conversation_id = ?", conversationID).
		Find(&messages).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询浏览器智能体消息列表失败")
	}
	return
}

func (r *BrowserAgentDB) ListMessagesIDListByConversationID(ctx context.Context, conversationID int64) (messageIDList []int64, err error) {
	if err = DB(ctx, r.db).
		Model(&entity.BrowserAgentMessage{}).
		Where("conversation_id = ?", conversationID).
		Pluck("id", &messageIDList).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询浏览器智能体消息ID列表失败")
	}
	return
}

// =========================
// Action CRUD
// =========================

func (r *BrowserAgentDB) CreateAction(ctx context.Context, action *entity.BrowserAgentAction) error {
	if err := DB(ctx, r.db).Create(action).Error; err != nil {
		return errors.WrapDBError(err, "创建浏览器智能体操作失败")
	}
	return nil
}

func (r *BrowserAgentDB) CreateActions(ctx context.Context, actions []*entity.BrowserAgentAction) error {
	if len(actions) == 0 {
		return nil
	}
	if err := DB(ctx, r.db).Create(&actions).Error; err != nil {
		return errors.WrapDBError(err, "批量创建浏览器智能体操作失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetActionByID(ctx context.Context, id int64) (action *entity.BrowserAgentAction, err error) {
	if err = DB(ctx, r.db).Where("id = ?", id).First(&action).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errors.NewDBError("操作不存在")
		}
		return nil, errors.WrapDBError(err, "查询浏览器智能体操作失败")
	}
	return
}

func (r *BrowserAgentDB) ListActionsByMessageID(ctx context.Context, messageID int64) (actions []*entity.BrowserAgentAction, err error) {
	if err = DB(ctx, r.db).Model(&entity.BrowserAgentAction{}).
		Where("message_id = ?", messageID).
		Find(&actions).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询浏览器智能体操作列表失败")
	}
	return
}

func (r *BrowserAgentDB) UpdateAction(ctx context.Context, action *entity.BrowserAgentAction) error {
	if err := DB(ctx, r.db).Save(action).Error; err != nil {
		return errors.WrapDBError(err, "更新浏览器智能体操作失败")
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

	if err := DB(ctx, r.db).Model(&entity.BrowserAgentAction{}).
		Where("id = ?", id).Updates(updates).Error; err != nil {
		return errors.WrapDBError(err, "更新操作状态失败")
	}
	return nil
}

func (r *BrowserAgentDB) GetPendingActionsByMessageID(ctx context.Context, messageID int64) (actions []*entity.BrowserAgentAction, err error) {
	if err = DB(ctx, r.db).Model(&entity.BrowserAgentAction{}).
		Where("message_id = ? AND status = ?", messageID, entity.ActionStatusPending).
		Order("sequence ASC").
		Find(&actions).Error; err != nil {
		return nil, errors.WrapDBError(err, "查询待执行操作失败")
	}
	return
}

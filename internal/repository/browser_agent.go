package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/repository/db"
	"context"
)

type BrowserAgentRepo struct {
	*db.BrowserAgentDB
}

func (r *BrowserAgentRepo) CreateConversation(ctx context.Context, conv *entity.BrowserAgentConversation) error {
	return r.BrowserAgentDB.CreateConversation(ctx, conv)
}

func (r *BrowserAgentRepo) GetConversationByID(ctx context.Context, id int64) (*entity.BrowserAgentConversation, error) {
	return r.BrowserAgentDB.GetConversationByID(ctx, id)
}

func (r *BrowserAgentRepo) ListConversationsByUserID(ctx context.Context, userID int64, queryParam *query.BrowserAgentConversation) ([]*entity.BrowserAgentConversation, int64, error) {
	return r.BrowserAgentDB.ListConversationsByUserID(ctx, userID, queryParam)
}

func (r *BrowserAgentRepo) UpdateConversation(ctx context.Context, conv *entity.BrowserAgentConversation) error {
	return r.BrowserAgentDB.UpdateConversation(ctx, conv)
}

func (r *BrowserAgentRepo) DeleteConversation(ctx context.Context, id int64) error {
	return r.BrowserAgentDB.DeleteConversation(ctx, id)
}

func (r *BrowserAgentRepo) UpdateConversationState(ctx context.Context, id int64, state string) error {
	return r.BrowserAgentDB.UpdateConversationState(ctx, id, state)
}

func (r *BrowserAgentRepo) CreateMessage(ctx context.Context, msg *entity.BrowserAgentMessage) error {
	return r.BrowserAgentDB.CreateMessage(ctx, msg)
}

func (r *BrowserAgentRepo) GetMessageByID(ctx context.Context, id int64) (*entity.BrowserAgentMessage, error) {
	return r.BrowserAgentDB.GetMessageByID(ctx, id)
}

func (r *BrowserAgentRepo) ListMessagesByConversationID(ctx context.Context, conversationID int64) ([]*entity.BrowserAgentMessage, error) {
	return r.BrowserAgentDB.ListMessagesByConversationID(ctx, conversationID)
}

func (r *BrowserAgentRepo) GetRecentMessages(ctx context.Context, conversationID int64, limit int) ([]*entity.BrowserAgentMessage, error) {
	return r.BrowserAgentDB.GetRecentMessages(ctx, conversationID, limit)
}

func (r *BrowserAgentRepo) CreateAction(ctx context.Context, action *entity.BrowserAgentAction) error {
	return r.BrowserAgentDB.CreateAction(ctx, action)
}

func (r *BrowserAgentRepo) CreateActions(ctx context.Context, actions []*entity.BrowserAgentAction) error {
	return r.BrowserAgentDB.CreateActions(ctx, actions)
}

func (r *BrowserAgentRepo) GetActionByID(ctx context.Context, id int64) (*entity.BrowserAgentAction, error) {
	return r.BrowserAgentDB.GetActionByID(ctx, id)
}

func (r *BrowserAgentRepo) ListActionsByMessageID(ctx context.Context, messageID int64) ([]*entity.BrowserAgentAction, error) {
	return r.BrowserAgentDB.ListActionsByMessageID(ctx, messageID)
}

func (r *BrowserAgentRepo) UpdateAction(ctx context.Context, action *entity.BrowserAgentAction) error {
	return r.BrowserAgentDB.UpdateAction(ctx, action)
}

func (r *BrowserAgentRepo) UpdateActionStatus(ctx context.Context, id int64, status string, errMsg *string, execTime *int) error {
	return r.BrowserAgentDB.UpdateActionStatus(ctx, id, status, errMsg, execTime)
}

func (r *BrowserAgentRepo) GetPendingActionsByMessageID(ctx context.Context, messageID int64) ([]*entity.BrowserAgentAction, error) {
	return r.BrowserAgentDB.GetPendingActionsByMessageID(ctx, messageID)
}

func (r *BrowserAgentRepo) GetMessageWithActions(ctx context.Context, messageID int64) (*entity.BrowserAgentMessage, []*entity.BrowserAgentAction, error) {
	return r.BrowserAgentDB.GetMessageWithActions(ctx, messageID)
}

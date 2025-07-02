package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type AgentFileDB struct {
	db *gorm.DB
}

func NewAgentFileDB(db *gorm.DB) *AgentFileDB {
	return &AgentFileDB{
		db: db,
	}
}

func (a *AgentFileDB) CreateAgentFile(c context.Context, e *entity.AgentFile) (err error) {
	if err = DB(c, a.db).Create(e).Error; err != nil {
		err = errors.WrapDBError(err, "创建AI知识库文件失败")
		return
	}
	return
}

func (a *AgentFileDB) GetAgentFileIDsByAgentID(c context.Context, agentId int64) (ids []int64, err error) {
	if err = DB(c, a.db).
		Model(&entity.AgentFile{}).
		Select("id").
		Where("agent_id = ?", agentId).
		Scan(&ids).Error; err != nil {
		err = errors.WrapDBError(err, "获取AI知识库文件ID失败")
		return
	}
	return
}

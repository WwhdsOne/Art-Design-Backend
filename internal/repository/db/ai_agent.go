package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type AIAgentDB struct {
	db *gorm.DB
}

func NewAIAgentDB(db *gorm.DB) *AIAgentDB {
	return &AIAgentDB{
		db: db,
	}
}

func (a *AIAgentDB) Create(c context.Context, e *entity.AIAgent) (err error) {
	if err = DB(c, a.db).Create(e).Error; err != nil {
		err = errors.WrapDBError(err, "创建AI模型失败")
		return
	}
	return
}

func (a *AIAgentDB) GetAIAgentPage(c context.Context, q *query.AIAgent) (pageRes []*entity.AIAgent, total int64, err error) {
	db := DB(c, a.db)

	// 构建通用查询条件
	queryConditions := db.Model(&entity.AIAgent{})

	if q.Name != nil {
		queryConditions = queryConditions.Where("name LIKE ?", "%"+*q.Name+"%")
	}

	// 查询总数
	if err = queryConditions.Count(&total).Error; err != nil {
		err = errors.WrapDBError(err, "获取智能体分页失败")
		return
	}

	// 查询分页数据
	if err = queryConditions.Scopes(q.Paginate()).Find(&pageRes).Error; err != nil {
		err = errors.WrapDBError(err, "获取智能体分页失败")
		return
	}
	return
}

func (a *AIAgentDB) GetSimpleAgentList(c context.Context) (agentList []*entity.AIAgent, err error) {
	if err = DB(c, a.db).Select("id", "name").Find(&agentList).Error; err != nil {
		err = errors.WrapDBError(err, "获取智能体列表失败")
		return
	}
	return
}

func (a *AIAgentDB) GetAgentByID(c context.Context, id int64) (agent *entity.AIAgent, err error) {
	if err = DB(c, a.db).Where("id = ?", id).First(&agent).Error; err != nil {
		err = errors.WrapDBError(err, "获取智能体失败")
		return
	}
	return
}

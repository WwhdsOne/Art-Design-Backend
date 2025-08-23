package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/errors"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type AIModelDB struct {
	db *gorm.DB
}

func NewAIModelDB(db *gorm.DB) *AIModelDB {
	return &AIModelDB{
		db: db,
	}
}

func (a *AIModelDB) CheckAIDuplicate(c context.Context, model *entity.AIModel) (err error) {
	var result struct {
		ModelExists   bool
		BaseURLExists bool
		ModelIDExists bool
	}

	// 假设 AIModel 使用字符串 ID，可选排除逻辑
	excludeID := ""
	if model.ID != 0 {
		excludeID = fmt.Sprintf("AND id != '%d'", model.ID)
	}

	var queryCondition strings.Builder
	args := make([]interface{}, 0)
	conditions := make([]string, 0)

	if model.Model != "" {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM ai_model WHERE model = ? "+excludeID+") AS model_exists")
		args = append(args, model.Model)
	}
	if model.ModelID != "" {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM ai_model WHERE model_id = ? "+excludeID+") AS model_id_exists")
		args = append(args, model.ModelID)
	}

	if len(conditions) == 0 {
		return
	}

	queryCondition.WriteString("SELECT ")
	queryCondition.WriteString(strings.Join(conditions, ", "))

	if err = DB(c, a.db).Raw(queryCondition.String(), args...).Scan(&result).Error; err != nil {
		return
	}

	switch {
	case result.ModelExists:
		err = errors.NewDBError("模型名称重复")
	case result.ModelIDExists:
		err = errors.NewDBError("模型接口标识重复")
	}
	return
}

func (a *AIModelDB) Create(c context.Context, e *entity.AIModel) (err error) {
	if err = DB(c, a.db).Create(e).Error; err != nil {
		err = errors.WrapDBError(err, "创建AI模型失败")
		return
	}
	return
}

func (a *AIModelDB) GetAIModelByID(c context.Context, id int64) (res *entity.AIModel, err error) {
	if err = DB(c, a.db).Where("id = ?", id).First(&res).Error; err != nil {
		err = errors.WrapDBError(err, "查询AI模型失败")
		return
	}
	return
}

func (a *AIModelDB) GetAIModelPage(c context.Context, q *query.AIModel) (pageRes []*entity.AIModel, total int64, err error) {
	db := DB(c, a.db)

	// 构建通用查询条件
	queryConditions := db.Model(&entity.AIModel{})

	if q.Model != nil {
		queryConditions = queryConditions.Where("model LIKE ?", "%"+*q.Model+"%")
	}
	if q.ModelType != nil {
		queryConditions = queryConditions.Where("model_type LIKE ?", *q.ModelType)
	}
	if q.Provider != nil {
		queryConditions = queryConditions.Where("provider LIKE ?", *q.Provider)
	}
	if q.Enabled != nil {
		queryConditions = queryConditions.Where("enabled = ?", *q.Enabled)
	}

	// 查询总数
	if err = queryConditions.Count(&total).Error; err != nil {
		err = errors.WrapDBError(err, "获取模型分页失败")
		return
	}

	// 查询分页数据
	if err = queryConditions.Scopes(q.Paginate()).Find(&pageRes).Error; err != nil {
		err = errors.WrapDBError(err, "获取模型分页失败")
		return
	}
	return
}

func (a *AIModelDB) GetSimpleChatModelList(c context.Context) (models []*entity.AIModel, err error) {
	if err = DB(c, a.db).Select("id", "icon", "model").
		Where("enabled = ?", true).
		Where("model_type = ?", "chat").Find(&models).Error; err != nil {
		err = errors.WrapDBError(err, "获取模型简洁列表失败")
		return
	}
	return
}

func (a *AIModelDB) GetRerankModel(c context.Context) (model *entity.AIModel, err error) {
	if err = DB(c, a.db).Where("enabled = ?", true).Where("model_type = ?", "rerank").First(&model).Error; err != nil {
		err = errors.WrapDBError(err, "获取模型失败")
		return
	}
	return
}

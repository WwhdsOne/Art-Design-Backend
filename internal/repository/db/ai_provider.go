package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/errors"
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type AIProviderDB struct {
	db *gorm.DB
}

func NewAIProviderDB(db *gorm.DB) *AIProviderDB {
	return &AIProviderDB{
		db: db,
	}
}

func (a *AIProviderDB) CheckAIDuplicate(c context.Context, provider *entity.AIProvider) (err error) {
	var result struct {
		NameExists bool
	}

	// 假设 AIProvider 使用字符串 ID，可选排除逻辑
	excludeID := ""
	if provider.ID != 0 {
		excludeID = fmt.Sprintf("AND id != '%d'", provider.ID)
	}

	var queryCondition strings.Builder
	args := make([]interface{}, 0)
	conditions := make([]string, 0)

	if provider.Name != "" {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM ai_model WHERE name = ? "+excludeID+") AS name_exists")
		args = append(args, provider.Name)
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
	case result.NameExists:
		err = errors.NewDBError("模型名称重复")
	}
	return
}

func (a *AIProviderDB) Create(c context.Context, provider *entity.AIProvider) (err error) {
	if err = DB(c, a.db).Create(provider).Error; err != nil {
		err = errors.WrapDBError(err, "创建AI供应商失败")
		return
	}
	return
}

func (a *AIProviderDB) GetAIProviderPage(c context.Context, q *query.AIProvider) (res []*entity.AIProvider, total int64, err error) {
	db := DB(c, a.db)

	// 构建通用查询条件
	queryConditions := db.Model(&entity.AIProvider{})
	if err = DB(c, a.db).Model(&entity.AIProvider{}).Count(&total).Error; err != nil {
		err = errors.WrapDBError(err, "获取AI供应商分页失败")
		return
	}
	if q.Name != nil {
		queryConditions.Where("name LIKE ?", "%"+*q.Name+"%")
	}
	if q.Enabled != nil {
		queryConditions.Where("enabled = ?", *q.Enabled)
	}

	if err = queryConditions.Scopes(q.Paginate()).Find(&res).Error; err != nil {
		err = errors.WrapDBError(err, "获取AI供应商分页失败")
		return
	}
	return
}

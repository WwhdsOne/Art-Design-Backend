package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/repository/dupcheck"
	"Art-Design-Backend/pkg/errors"
	"context"

	"gorm.io/gorm"
)

type AIProviderDB struct {
	db *gorm.DB
}

func NewAIProviderDB(db *gorm.DB) *AIProviderDB {
	return &AIProviderDB{
		db: db,
	}
}

func (a *AIProviderDB) CheckAIDuplicate(c context.Context, provider *entity.AIProvider) error {
	return dupcheck.Check(DB(c, a.db), "ai_provider", provider.ID, []dupcheck.Field{
		{Column: "name", ErrMsg: "模型名称重复", Value: provider.Name},
	})
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

func (a *AIProviderDB) GetSimpleProviderList(c context.Context) (providerList []*entity.AIProvider, err error) {
	if err = DB(c, a.db).Select("id", "name").Find(&providerList).Error; err != nil {
		err = errors.WrapDBError(err, "获取AI供应商列表失败")
		return
	}
	return
}

func (a *AIProviderDB) GetProviderNameByIDList(c context.Context, idList []int64) (providerList []*entity.AIProvider, err error) {
	if err = DB(c, a.db).Where("id in ?", idList).Find(&providerList).Error; err != nil {
		err = errors.WrapDBError(err, "获取AI供应商名称列表失败")
		return
	}
	return
}

func (a *AIProviderDB) GetProviderByID(c context.Context, id int64) (provider *entity.AIProvider, err error) {
	if err = DB(c, a.db).Where("id = ?", id).Where("enabled = ?", true).First(&provider).Error; err != nil {
		err = errors.WrapDBError(err, "获取AI供应商失败")
		return
	}
	return
}

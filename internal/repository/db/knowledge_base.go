package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/errors"
	"context"

	"gorm.io/gorm"
)

type KnowledgeBaseDB struct {
	db *gorm.DB
}

func NewKnowledgeBaseDB(db *gorm.DB) *KnowledgeBaseDB {
	return &KnowledgeBaseDB{
		db: db,
	}
}

func (k *KnowledgeBaseDB) CreateKnowledgeFile(c context.Context, e *entity.KnowledgeBaseFile) (err error) {
	if err = DB(c, k.db).Create(e).Error; err != nil {
		err = errors.WrapDBError(err, "创建知识库文件失败")
		return
	}
	return
}

func (k *KnowledgeBaseDB) GetKnowledgeFilePage(c context.Context, req *query.KnowledgeBaseFile, CreateUserIDList []int64) (res []*entity.KnowledgeBaseFile, total int64, err error) {
	db := DB(c, k.db)

	queryCondition := db.Model(&entity.KnowledgeBaseFile{})
	if req.FileType != nil {
		queryCondition = queryCondition.Where("file_type like ?", req.FileType)
	}
	if req.OriginalFileName != nil {
		queryCondition = queryCondition.Where("original_file_name like ?", "%"+*req.OriginalFileName+"%")
	}
	if len(CreateUserIDList) != 0 {
		queryCondition = queryCondition.Where("created_by in ?", CreateUserIDList)
	}
	// 前端返回的是时间数组，第一个是开始时间，第二个是结束时间
	if req.TimeRange != nil {
		queryCondition = queryCondition.Where("created_at between ? and ?", req.TimeRange[0], req.TimeRange[1])
	}
	if err = queryCondition.Count(&total).Error; err != nil {
		err = errors.WrapDBError(err, "获取知识库文件数量失败")
		return
	}
	if err = queryCondition.Scopes(req.Paginate()).Find(&res).Error; err != nil {
		err = errors.WrapDBError(err, "获取知识库文件失败")
		return
	}

	return
}

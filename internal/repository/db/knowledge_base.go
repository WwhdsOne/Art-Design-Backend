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

func (k *KnowledgeBaseDB) GetKnowledgeBasePage(c context.Context, q *query.KnowledgeBase, createUser int64) (res []*entity.KnowledgeBase, err error) {
	if err = DB(c, k.db).Scopes(q.Paginate()).
		Where("created_by = ?", createUser).
		Find(&res).Error; err != nil {
		err = errors.WrapDBError(err, "获取知识库失败")
		return
	}
	return
}

func (k *KnowledgeBaseDB) CreateKnowledgeBase(ctx context.Context, e *entity.KnowledgeBase) error {
	if err := DB(ctx, k.db).Create(e).Error; err != nil {
		return errors.WrapDBError(err, "创建知识库失败")
	}
	return nil
}

func (k *KnowledgeBaseDB) DeleteKnowledgeBase(ctx context.Context, id int64) error {
	if err := DB(ctx, k.db).Where("id = ?", id).Delete(&entity.KnowledgeBase{}).Error; err != nil {
		return errors.WrapDBError(err, "删除知识库失败")
	}
	return nil
}

func (k *KnowledgeBaseDB) UpdateKnowledgeBase(c context.Context, e *entity.KnowledgeBase) error {
	if err := DB(c, k.db).Where("id = ?", e.ID).Updates(e).Error; err != nil {
		return errors.WrapDBError(err, "更新知识库失败")
	}
	return nil
}

func (k *KnowledgeBaseDB) GetKnowledgeBaseFileByIDs(c context.Context, ids []int64) (res []*entity.KnowledgeBaseFile, err error) {
	if err = DB(c, k.db).Where("id in ?", ids).Find(&res).Error; err != nil {
		err = errors.WrapDBError(err, "获取知识库文件失败")
		return
	}
	return
}

func (k *KnowledgeBaseDB) GetSimpleKnowledgeBaseList(c context.Context, userID int64) (res []*entity.KnowledgeBase, err error) {
	if err = DB(c, k.db).Select("id", "name").
		Where("created_by = ?", userID).Find(&res).Error; err != nil {
		err = errors.WrapDBError(err, "获取知识库列表失败")
		return
	}
	return
}

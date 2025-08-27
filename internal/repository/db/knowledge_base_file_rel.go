package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"

	"gorm.io/gorm"
)

type KnowledgeBaseFileRelDB struct {
	db *gorm.DB
}

func NewKnowledgeBaseFileRelDB(db *gorm.DB) *KnowledgeBaseFileRelDB {
	return &KnowledgeBaseFileRelDB{
		db: db,
	}
}
func (k *KnowledgeBaseFileRelDB) CreateKnowledgeBaseFileRel(c context.Context, rel []*entity.KnowledgeBaseFileRel) (err error) {
	if err = DB(c, k.db).Create(rel).Error; err != nil {
		return errors.NewDBError("创建知识库文件关系失败")
	}
	return
}

func (k *KnowledgeBaseFileRelDB) DeleteKnowledgeBaseFileRel(c context.Context, id int64) (err error) {
	if err = DB(c, k.db).Where("knowledge_base_id = ?", id).Delete(&entity.KnowledgeBaseFileRel{}).Error; err != nil {
		return errors.NewDBError("删除知识库文件关系失败")
	}
	return
}

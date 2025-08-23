package db

import (
	"Art-Design-Backend/internal/model/entity"
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

func (k *KnowledgeBaseDB) CreateAgentFile(c context.Context, e *entity.KnowledgeBaseFile) (err error) {
	if err = DB(c, k.db).Create(e).Error; err != nil {
		err = errors.WrapDBError(err, "创建知识库文件失败")
		return
	}
	return
}

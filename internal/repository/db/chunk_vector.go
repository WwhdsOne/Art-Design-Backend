package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type ChunkVectorDB struct {
	db *gorm.DB
}

func NewChunkVectorDB(db *gorm.DB) *ChunkVectorDB {
	return &ChunkVectorDB{
		db: db,
	}
}

func (f *ChunkVectorDB) CreateChunkVector(c context.Context, e *entity.ChunkVector) (err error) {
	if err = DB(c, f.db).Create(e).Error; err != nil {
		err = errors.WrapDBError(err, "创建文件块向量失败")
		return
	}
	return
}

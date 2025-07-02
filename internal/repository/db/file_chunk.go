package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type FileChunkDB struct {
	db *gorm.DB
}

func NewFileChunkDB(db *gorm.DB) *FileChunkDB {
	return &FileChunkDB{
		db: db,
	}
}

func (f *FileChunkDB) CreateFileChunk(c context.Context, e *entity.FileChunk) (err error) {
	if err = DB(c, f.db).Create(e).Error; err != nil {
		err = errors.WrapDBError(err, "创建文件块失败")
		return
	}
	return
}

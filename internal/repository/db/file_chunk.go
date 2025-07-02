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

func (f *FileChunkDB) GetFileChunkIDsByFileIDList(c context.Context, fileIDList []int64) (res []int64, err error) {
	if err = DB(c, f.db).
		Model(&entity.FileChunk{}).
		Select("id").
		Where("file_id IN (?)", fileIDList).
		Scan(&res).Error; err != nil {
		err = errors.WrapDBError(err, "获取文件块ID失败")
		return
	}
	return
}

func (f *FileChunkDB) GetFileContentByIDList(c context.Context, ids []int64) (contents []string, err error) {
	if err = DB(c, f.db).
		Model(&entity.FileChunk{}). // 指定表
		Select("content").          // 只查询 content 字段
		Where("id IN (?)", ids).
		Pluck("content", &contents). // 直接提取 content 字段到字符串切片
		Error; err != nil {
		err = errors.WrapDBError(err, "获取文件块内容失败")
		return
	}
	return
}

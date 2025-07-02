package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/errors"
	"context"
	"fmt"
	"github.com/pgvector/pgvector-go"
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

func (f *ChunkVectorDB) SearchTopKByEmbedding(
	c context.Context,
	chunkIDList []int64, // 待筛选的 chunk ID 列表，只在这些 chunk 中查找最相似的向量
	queryVec pgvector.Vector, // 查询向量，用于计算相似度
	topK int, // 返回的最相似向量数量限制
) (similarChunkIDs []int64, err error) {

	// 如果传入的 chunkIDList 为空，直接返回空结果
	if len(chunkIDList) == 0 {
		return nil, nil
	}

	// 使用原生 SQL 查询
	err = DB(c, f.db).
		Raw(`
			SELECT chunk_id
			FROM chunk_vectors
			WHERE chunk_id IN (?)          -- 只在指定的 chunkIDList 范围内搜索
			ORDER BY embedding <#> ? -- 以内积（Inner Product）相似度排序，越大越相似
			LIMIT ?                 -- 限制返回 topK 条结果
		`, chunkIDList, queryVec, topK).
		Scan(&similarChunkIDs).Error // 扫描结果到 ID 列表

	if err != nil {
		return nil, fmt.Errorf("向量搜索失败: %w", err)
	}

	return similarChunkIDs, nil
}

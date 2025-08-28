package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/db"
	"context"

	"github.com/pgvector/pgvector-go"
)

type KnowledgeBaseRepo struct {
	*db.KnowledgeBaseDB
	*db.FileChunkDB
	*db.ChunkVectorDB
	*db.KnowledgeBaseFileRelDB
}

func (k *KnowledgeBaseRepo) GetKnowledgeBaseFilesByID(ctx context.Context, id int64) (files []*entity.KnowledgeBaseFile, err error) {
	var knowledgeBaseFileRelList []*entity.KnowledgeBaseFileRel
	knowledgeBaseFileRelList, err = k.KnowledgeBaseFileRelDB.GetKnowledgeBaseFileRel(ctx, id)
	if err != nil {
		return
	}
	knowledgeBaseFileIDs := make([]int64, 0)
	for _, knowledgeBaseFileRel := range knowledgeBaseFileRelList {
		knowledgeBaseFileIDs = append(knowledgeBaseFileIDs, knowledgeBaseFileRel.KnowledgeBaseFileID)
	}
	if files, err = k.GetKnowledgeBaseFileByIDs(ctx, knowledgeBaseFileIDs); err != nil {
		return
	}
	return
}

func (k *KnowledgeBaseRepo) SearchAgentRelatedChunks(c context.Context, knowledgeBaseID int64, vector []float32) (chunks []string, err error) {
	// 获取智能体的文件ID列表
	knowledgeBaseFileIDList, err := k.GetKnowledgeBaseFilesByID(c, knowledgeBaseID)
	if err != nil {
		return
	}
	knowledgeBaseFileIDs := make([]int64, 0, len(knowledgeBaseFileIDList))
	for _, knowledgeBaseFile := range knowledgeBaseFileIDList {
		knowledgeBaseFileIDs = append(knowledgeBaseFileIDs, knowledgeBaseFile.ID)
	}
	// 获取文件块
	fileChunkIDs, err := k.FileChunkDB.GetFileChunkIDsByFileIDList(c, knowledgeBaseFileIDs)
	if err != nil {
		return
	}
	// 根据文件块ID和向量查询出对应的ID
	similarChunkIDs, err := k.ChunkVectorDB.SearchTopKByEmbedding(c, fileChunkIDs, pgvector.NewVector(vector), 10)
	if err != nil {
		return
	}
	// 查询文件块内容
	chunks, err = k.FileChunkDB.GetFileContentByIDList(c, similarChunkIDs)
	if err != nil {
		return
	}
	return
}

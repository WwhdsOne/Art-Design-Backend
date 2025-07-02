package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/db"
	"context"
	"github.com/pgvector/pgvector-go"
)

type AIAgentRepo struct {
	*db.AIAgentDB
	*db.AgentFileDB
	*db.FileChunkDB
	*db.ChunkVectorDB
}

func (a *AIAgentRepo) GetAIAgentByIDWithCache(c context.Context, id int64) (res *entity.AIAgent, err error) {
	res, err = a.AIAgentDB.GetAgentByID(c, id)
	return res, err
}

func (a *AIAgentRepo) SearchAgentRelatedChunks(c context.Context, agentID int64, vector []float32) (chunks []string, err error) {
	// 获取智能体的文件ID列表
	agentFileIDList, err := a.GetAgentFileIDsByAgentID(c, agentID)
	if err != nil {
		return
	}
	// 获取文件块
	fileChunkIDs, err := a.GetFileChunkIDsByFileIDList(c, agentFileIDList)
	if err != nil {
		return
	}
	// 根据文件块ID和向量查询出对应的ID
	similarChunkIDs, err := a.SearchTopKByEmbedding(c, fileChunkIDs, pgvector.NewVector(vector), 3)
	if err != nil {
		return
	}
	// 查询文件块内容
	chunks, err = a.GetFileContentByIDList(c, similarChunkIDs)
	if err != nil {
		return
	}
	return
}

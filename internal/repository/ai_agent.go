package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"context"
	"errors"
	"github.com/pgvector/pgvector-go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AIAgentRepo struct {
	*db.AIAgentDB
	*db.AgentFileDB
	*db.FileChunkDB
	*db.ChunkVectorDB
	*cache.AIAgentCache
}

func (a *AIAgentRepo) GetAIAgentByIDWithCache(c context.Context, id int64) (res *entity.AIAgent, err error) {
	res, err = a.AIAgentCache.GetAIAgentInfo(id)
	if err == nil {
		// 缓存命中，直接返回
		return
	}
	if !errors.Is(err, redis.Nil) {
		// 缓存出错，但不是未命中，记录日志
		zap.L().Warn("获取AI模型缓存失败", zap.Error(err))
	}

	// 缓存未命中，查数据库
	res, err = a.AIAgentDB.GetAgentByID(c, id)
	if err != nil {
		return nil, err
	}

	// 异步回填缓存
	go func(agent *entity.AIAgent) {
		if agent == nil {
			return
		}
		if cacheErr := a.AIAgentCache.SetAgentInfo(agent); cacheErr != nil {
			zap.L().Warn("设置AI模型缓存失败", zap.Error(cacheErr))
		}
	}(res)

	return res, nil
}

func (a *AIAgentRepo) SearchAgentRelatedChunks(c context.Context, agentID int64, vector []float32) (chunks []string, err error) {
	// 获取智能体的文件ID列表
	agentFileIDList, err := a.AgentFileDB.GetAgentFileIDsByAgentID(c, agentID)
	if err != nil {
		return
	}
	// 获取文件块
	fileChunkIDs, err := a.FileChunkDB.GetFileChunkIDsByFileIDList(c, agentFileIDList)
	if err != nil {
		return
	}
	// 根据文件块ID和向量查询出对应的ID
	similarChunkIDs, err := a.ChunkVectorDB.SearchTopKByEmbedding(c, fileChunkIDs, pgvector.NewVector(vector), 3)
	if err != nil {
		return
	}
	// 查询文件块内容
	chunks, err = a.FileChunkDB.GetFileContentByIDList(c, similarChunkIDs)
	if err != nil {
		return
	}
	return
}

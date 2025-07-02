package repository

import (
	"Art-Design-Backend/internal/repository/db"
)

type AIAgentRepo struct {
	*db.AIAgentDB
	*db.AgentFileDB
	*db.FileChunkDB
	*db.ChunkVectorDB
}

package repository

import (
	"Art-Design-Backend/internal/repository/db"
)

type KnowledgeBaseRepo struct {
	*db.KnowledgeBaseDB
	*db.FileChunkDB
	*db.ChunkVectorDB
	*db.KnowledgeBaseFileRelDB
}

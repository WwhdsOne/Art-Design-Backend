package entity

import (
	"github.com/pgvector/pgvector-go"
)

// ChunkVector 分块向量模型
type ChunkVector struct {
	ID      int64           `gorm:"primaryKey;"`
	ChunkID int64           `gorm:"not null;index"`
	Vector  pgvector.Vector `gorm:"type:vector(1024);not null"` // 这里2560维根据你的模型调整
}

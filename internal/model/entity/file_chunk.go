package entity

type FileChunk struct {
	ID         int64  `gorm:"primaryKey"`
	FileID     int64  `gorm:"not null;index"`     // 外键，关联 agent_files 表
	ChunkIndex int    `gorm:"not null"`           // 分块索引
	Content    string `gorm:"type:text;not null"` // 分块内容
}

func (FileChunk) TableName() string {
	return "file_chunks"
}

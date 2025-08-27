package entity

import "Art-Design-Backend/pkg/constant/tablename"

type KnowledgeBaseFileRel struct {
	KnowledgeBaseID     int64 `gorm:"knowledge_base_id"`
	KnowledgeBaseFileID int64 `gorm:"knowledge_base_file_id"`
}

func (k *KnowledgeBaseFileRel) TableName() string {
	return tablename.KnowledgeBaseFileRelTableName
}

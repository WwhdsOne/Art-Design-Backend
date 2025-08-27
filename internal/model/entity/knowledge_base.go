package entity

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/pkg/constant/tablename"
)

type KnowledgeBase struct {
	common.BaseModel
	Name        string               `gorm:"type:varchar(50);not null;comment:知识库名称"`
	Description string               `gorm:"type:varchar(256);comment:备注"`
	Files       []*KnowledgeBaseFile `gorm:"many2many:knowledge_base_file_rel;comment:关联文件"`
}

func (k *KnowledgeBase) TableName() string {
	return tablename.KnowledgeBaseTableName
}

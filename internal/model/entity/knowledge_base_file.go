package entity

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/pkg/constant/tablename"
)

type KnowledgeBaseFile struct {
	common.BaseModel

	OriginalFileName string `gorm:"size:500;not null;comment:原始文件名"`                      // 用户上传时的文件名
	FileType         string `gorm:"type:varchar(20);not null;comment:文件类型(pdf/docx/txt)"` // 文件类型
	FileSize         int64  `gorm:"not null;comment:文件大小(字节)"`                            // 单位：字节
	FilePath         string `gorm:"size:256;not null;comment:存储路径"`                       // 文件存储路径
}

func (k *KnowledgeBaseFile) TableName() string {
	return tablename.KnowledgeBaseFileTableName
}

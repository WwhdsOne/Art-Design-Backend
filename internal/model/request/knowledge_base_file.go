package request

import "Art-Design-Backend/internal/model/common"

type KnowledgeBaseFile struct {
	ID               common.LongStringID `json:"id"`
	OriginalFileName string              `json:"original_file_name"`
	FileType         string              `json:"file_type"`
	FileSize         int64               `json:"file_size"`
	FilePath         string              `json:"file_path"`
}

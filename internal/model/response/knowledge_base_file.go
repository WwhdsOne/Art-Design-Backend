package response

import (
	"github.com/dromara/carbon/v2"
)

type KnowledgeBaseFile struct {
	ID               int64           `json:"id,string"`         // 主键ID
	OriginalFileName string          `json:"original_filename"` // 用户上传时的文件名
	FileType         string          `json:"file_type"`         // 文件类型(pdf/docx/txt)
	FileSize         int64           `json:"file_size"`         // 文件大小(字节)
	FilePath         string          `json:"file_path"`         // 存储路径
	CreateUser       string          `json:"create_user"`       // 创建人字段，记录创建操作者的标识
	CreatedAt        carbon.DateTime `json:"created_at"`        // 创建时间字段，记录记录创建时的时间戳
}

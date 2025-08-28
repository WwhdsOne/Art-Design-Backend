package request

import "Art-Design-Backend/internal/model/common"

type KnowledgeBase struct {
	ID          common.LongStringID  `json:"id" label:"知识库ID"`
	Name        string               `json:"name" binding:"required,min=1,max=100" label:"知识库名称"`
	Description string               `json:"description" binding:"omitempty,max=500" label:"知识库描述"`
	Files       []*KnowledgeBaseFile `json:"files" binding:"omitempty,dive" label:"关联文件"`
}

package base

// PaginationResp 泛型结构体，T 表示 Data 字段的数据类型
type PaginationResp[T any] struct {
	Size  int `json:"size"`
	Page  int `json:"page"`
	Data  []T `json:"data"`
	Total int `json:"total"`
}

func BuildPageResp[T any](data []T, total int64, pageReq PaginationQ) PaginationResp[T] {
	return PaginationResp[T]{
		Page:  pageReq.Page,
		Size:  pageReq.Size,
		Total: int(total),
		Data:  data,
	}
}

package base

type PaginationResp[T any] struct {
	Size  int `json:"size"`
	Page  int `json:"page"`
	Data  []T `json:"data"`
	Total int `json:"total"`
}

func BuildPageResp[T any](data []T, total int64, pageReq PaginationReq) *PaginationResp[T] {
	return &PaginationResp[T]{
		Page:  pageReq.Page,
		Size:  pageReq.Size,
		Total: int(total),
		Data:  data,
	}
}

package base

import "gorm.io/gorm"

type PaginationReq struct {
	Size int `json:"size"`
	Page int `json:"page"`
}

func (r *PaginationReq) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.Page <= 0 {
			r.Page = 1
		}
		switch {
		case r.Size <= 0:
			r.Size = 10
		}
		offset := (r.Page - 1) * r.Size
		return db.Offset(offset).Limit(r.Size)
	}
}

type PaginationResp[T any] struct {
	Size  int `json:"size"`
	Page  int `json:"page"`
	Data  []T `json:"data"`
	Total int `json:"total"`
}

func BuildPageResp[T any](data []T, total int64, pageReq PaginationReq) PaginationResp[T] {
	return PaginationResp[T]{
		Page:  pageReq.Page,
		Size:  pageReq.Size,
		Total: int(total),
		Data:  data,
	}
}

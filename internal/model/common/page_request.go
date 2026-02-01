package common

import "gorm.io/gorm"

type PaginationReq struct {
	Size int `json:"size"`
	Page int `json:"page"`
}

func (r *PaginationReq) Paginate() func(db *gorm.DB) *gorm.DB {
	// 分页
	return func(db *gorm.DB) *gorm.DB {
		// 默认第一页
		if r.Page <= 0 {
			r.Page = 1
		}
		// 默认每页10条
		if r.Size <= 0 {
			r.Size = 10
		}
		offset := (r.Page - 1) * r.Size
		return db.Offset(offset).Limit(r.Size)
	}
}

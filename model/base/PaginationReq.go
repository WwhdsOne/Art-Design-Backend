package base

import "gorm.io/gorm"

type PaginationQ struct {
	Size int `form:"size" json:"size"`
	Page int `form:"page" json:"page"`
}

func (r *PaginationQ) Paginate() func(db *gorm.DB) *gorm.DB {
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

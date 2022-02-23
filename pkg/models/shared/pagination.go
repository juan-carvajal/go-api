package shared

import "gorm.io/gorm"

type Pagination struct {
	Total    int `json:"total"`
	PageSize int `json:"page_size"`
	Offset   int `json:"offset"`
}

func Paginate(p Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.Offset).Limit(p.PageSize)
	}
}

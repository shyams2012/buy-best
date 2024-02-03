package Pagination

import (
	"fmt"

	"gorm.io/gorm"
)

func Paginate(page *int, limit *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page != nil && limit != nil {
			offset := *page
			pageSize := *limit
			offset_cal := (offset - 1) * pageSize
			return db.Offset(offset_cal).Limit(pageSize)
		} else if limit != nil {
			pageSize := *limit
			return db.Limit(pageSize)
		} else {
			return db
		}

	}
}

func test5() {
	fmt.Print("hello")
}

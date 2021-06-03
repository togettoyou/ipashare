package model

import (
	"gorm.io/gorm"
)

// PaginationQ 分页结构体
type PaginationQ struct {
	//每页显示的数量
	PageSize int `json:"page_size"`
	//当前页码
	Page int `json:"page"`
	//分页的数据内容
	Data interface{} `json:"data"`
	//全部的页码数量
	Total int64 `json:"total"`
}

// Count 通用计数
func Count(total *int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Count(total)
	}
}

// Paginate 通用分页
func Paginate(page, pageSize *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case *pageSize > 100:
			*pageSize = 100
		case *pageSize <= 0:
			*pageSize = 10
		}
		if *page < 1 {
			*page = 1
		}
		offset := (*page - 1) * (*pageSize)
		return db.Offset(offset).Limit(*pageSize)
	}
}

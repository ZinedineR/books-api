package pagination

import (
	"gorm.io/gorm"
)

func Paginate[T any](page, pageSize int, query *gorm.DB) (PaginationResult[T], error) {
	if pageSize < 0 {
		pageSize = -1
	}
	if page < 1 {
		page = 1
	}

	var total int64
	var data []*T
	if err := query.Model(new(T)).Count(&total).Error; err != nil {
		return PaginationResult[T]{}, err
	}
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&data).Error; err != nil {
		return PaginationResult[T]{}, err
	}

	var totalPage int64
	if pageSize > 0 {
		totalPage = total / int64(pageSize)
		if total%int64(pageSize) > 0 {
			totalPage++
		}
	} else if pageSize == 0 {
		totalPage = 0
	} else {
		totalPage = 1
	}

	return PaginationResult[T]{
		Page:             page,
		PageSize:         pageSize,
		TotalPage:        totalPage,
		TotalDataPerPage: int64(len(data)),
		TotalData:        total,
		Data:             data,
	}, nil
}

type PaginationResult[T any] struct {
	Page             int   // The current page
	PageSize         int   // The size of the page
	TotalPage        int64 // The total number of pages
	TotalDataPerPage int64 // The total number of data per page
	TotalData        int64 // The total number of data
	Data             []*T  // The actual data
}

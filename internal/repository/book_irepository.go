package repository

import (
	"books-api/internal/entity"
	"books-api/internal/model"
	"context"
	"gorm.io/gorm"
)

type BookRepository interface {
	// Example operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.Book) error
	FindByName(ctx context.Context, tx *gorm.DB, column, value string) (
		*entity.Book, error,
	)
	FindByPagination(
		ctx context.Context, tx *gorm.DB, page model.PaginationParam, order model.OrderParam,
		filter model.FilterParams,
	) (*model.PaginationData[entity.Book], error)
	FindByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Book, error)
	UpdateTx(ctx context.Context, tx *gorm.DB, data *entity.Book) error
	DeleteByIDTx(ctx context.Context, tx *gorm.DB, id string) error
}

package repository

import (
	"books-api/internal/entity"
	"books-api/internal/model"
	"context"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	// Example operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.Author) error
	FindByName(ctx context.Context, tx *gorm.DB, column, value string) (
		*entity.Author, error,
	)
	FindByPagination(
		ctx context.Context, tx *gorm.DB, page model.PaginationParam, order model.OrderParam,
		filter model.FilterParams,
	) (*model.PaginationData[entity.Author], error)
	FindByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Author, error)
	UpdateTx(ctx context.Context, tx *gorm.DB, data *entity.Author) error
	DeleteByIDTx(ctx context.Context, tx *gorm.DB, id string) error
}

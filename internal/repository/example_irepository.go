package repository

import (
	"boiler-plate-clean/internal/entity"
	"context"
	"gorm.io/gorm"
)

type ExampleRepository interface {
	// Example operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.Example) error
	FindByName(ctx context.Context, tx *gorm.DB, column, value string) (
		*entity.Example, error,
	)
	FindByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Example, error)
}

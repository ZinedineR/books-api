package repository

import (
	"books-api/internal/entity"
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	// Example operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.User) error
	FindByName(ctx context.Context, tx *gorm.DB, column, value string) (
		*entity.User, error,
	)
}

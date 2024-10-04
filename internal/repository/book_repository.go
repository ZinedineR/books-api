package repository

import (
	"books-api/internal/entity"
)

type BookSQLRepo struct {
	Repository[entity.Book]
}

func NewBookSQLRepository() BookRepository {
	return &BookSQLRepo{}
}

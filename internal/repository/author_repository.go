package repository

import (
	"books-api/internal/entity"
)

type AuthorSQLRepo struct {
	Repository[entity.Author]
}

func NewAuthorSQLRepository() AuthorRepository {
	return &AuthorSQLRepo{}
}

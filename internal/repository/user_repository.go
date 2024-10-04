package repository

import (
	"books-api/internal/entity"
)

type UserSQLRepo struct {
	Repository[entity.User]
}

func NewUserSQLRepository() UserRepository {
	return &UserSQLRepo{}
}

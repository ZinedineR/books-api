package repository

import (
	"boiler-plate-clean/internal/entity"
)

type ExampleSQLRepo struct {
	Repository[entity.Example]
}

func NewExampleSQLRepository() ExampleRepository {
	return &ExampleSQLRepo{}
}

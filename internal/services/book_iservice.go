package service

import (
	"books-api/internal/entity"
	"books-api/internal/model"
	"books-api/pkg/exception"
	"context"
)

type BookService interface {
	// CRUD operations for Book
	Create(
		ctx context.Context, model *entity.UpsertBook,
	) *exception.Exception
	Update(
		ctx context.Context, id string, model *entity.UpsertBook,
	) *exception.Exception
	Delete(
		ctx context.Context, id string,
	) *exception.Exception
	List(ctx context.Context, req model.ListReq) (
		*ListBookResp, *exception.Exception,
	)
	FindOne(ctx context.Context, id string) (*entity.Book, *exception.Exception)
}

type ListBookResp struct {
	Pagination *model.Pagination `json:"pagination"`
	Data       []*entity.Book    `json:"data"`
}

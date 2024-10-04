package service

import (
	"books-api/internal/entity"
	"books-api/internal/model"
	"books-api/pkg/exception"
	"context"
)

type AuthorService interface {
	// CRUD operations for Author
	Create(
		ctx context.Context, model *entity.UpsertAuthor,
	) *exception.Exception
	Update(
		ctx context.Context, id string, model *entity.UpsertAuthor,
	) *exception.Exception
	Delete(
		ctx context.Context, id string,
	) *exception.Exception
	List(ctx context.Context, req model.ListReq) (
		*ListAuthorResp, *exception.Exception,
	)
	FindOne(ctx context.Context, id string) (*entity.Author, *exception.Exception)
}

type ListAuthorResp struct {
	Pagination *model.Pagination `json:"pagination"`
	Data       []*entity.Author  `json:"data"`
}

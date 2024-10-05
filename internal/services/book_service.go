package service

import (
	"books-api/internal/entity"
	"books-api/internal/model"
	"books-api/internal/repository"
	"context"
	//"books-api/pkg/exception"
	"books-api/pkg/exception"
	"books-api/pkg/xvalidator"
	"gorm.io/gorm"
)

type BookServiceImpl struct {
	db       *gorm.DB
	bookRepo repository.BookRepository
	validate *xvalidator.Validator
}

func NewBookService(
	db *gorm.DB, repo repository.BookRepository,
	validate *xvalidator.Validator,
) BookService {
	return &BookServiceImpl{
		db:       db,
		bookRepo: repo,
		validate: validate,
	}
}

func (s *BookServiceImpl) Create(
	ctx context.Context, model *entity.UpsertBook,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(model); errs != nil {
		return exception.InvalidArgument(errs)
	}
	body := &entity.Book{}
	body.GenerateModel(model)
	if err := s.bookRepo.CreateTx(ctx, tx, body); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *BookServiceImpl) Update(
	ctx context.Context, id string, model *entity.UpsertBook,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(model); errs != nil {
		return exception.InvalidArgument(errs)
	}
	body := &entity.Book{Id: id}
	body.GenerateModel(model)
	if err := s.bookRepo.UpdateTx(ctx, tx, body); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *BookServiceImpl) Delete(
	ctx context.Context, id string,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.bookRepo.DeleteByIDTx(ctx, tx, id); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *BookServiceImpl) List(ctx context.Context, req model.ListReq) (
	*ListBookResp, *exception.Exception,
) {
	result, err := s.bookRepo.FindByPagination(ctx, s.db, req.Page, req.Order, req.Filter)
	if err != nil {
		return nil, exception.Internal("failed to get book", err)
	}
	return &ListBookResp{
		Pagination: &model.Pagination{
			Page:             result.Page,
			PageSize:         result.PageSize,
			TotalPage:        result.TotalPage,
			TotalDataPerPage: result.TotalDataPerPage,
			TotalData:        result.TotalData,
		},
		Data: result.Data,
	}, nil
}

func (s *BookServiceImpl) FindOne(ctx context.Context, id string) (*entity.Book, *exception.Exception) {
	result, err := s.bookRepo.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	return result, nil
}

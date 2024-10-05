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

type AuthorServiceImpl struct {
	db         *gorm.DB
	authorRepo repository.AuthorRepository
	validate   *xvalidator.Validator
}

func NewAuthorService(
	db *gorm.DB, repo repository.AuthorRepository,
	validate *xvalidator.Validator,
) AuthorService {
	return &AuthorServiceImpl{
		db:         db,
		authorRepo: repo,
		validate:   validate,
	}
}

func (s *AuthorServiceImpl) Create(
	ctx context.Context, model *entity.UpsertAuthor,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(model); errs != nil {
		return exception.InvalidArgument(errs)
	}
	body := &entity.Author{}
	body.GenerateModel(model)
	if err := s.authorRepo.CreateTx(ctx, tx, body); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *AuthorServiceImpl) Update(
	ctx context.Context, id string, model *entity.UpsertAuthor,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(model); errs != nil {
		return exception.InvalidArgument(errs)
	}
	body := &entity.Author{Id: id}
	body.GenerateModel(model)
	if err := s.authorRepo.UpdateTx(ctx, tx, body); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *AuthorServiceImpl) Delete(
	ctx context.Context, id string,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.authorRepo.DeleteByIDTx(ctx, tx, id); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *AuthorServiceImpl) List(ctx context.Context, req model.ListReq) (
	*ListAuthorResp, *exception.Exception,
) {
	result, err := s.authorRepo.FindByPagination(ctx, s.db, req.Page, req.Order, req.Filter)
	if err != nil {
		return nil, exception.Internal("failed to get author", err)
	}
	return &ListAuthorResp{
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

func (s *AuthorServiceImpl) FindOne(ctx context.Context, id string) (*entity.Author, *exception.Exception) {
	result, err := s.authorRepo.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	return result, nil
}

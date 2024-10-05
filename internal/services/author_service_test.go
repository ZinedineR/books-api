package service_test

import (
	"books-api/internal/entity"
	"books-api/internal/mocks"
	"books-api/internal/model"
	service "books-api/internal/services"
	"books-api/pkg/xvalidator"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateAuthor(t *testing.T) {
	mockAppCtx := context.Background()
	t.Run("CreateAuthor Success", func(t *testing.T) {
		// Set up input
		request := &entity.UpsertAuthor{
			Name:      "Ernest Hemingway",
			Birthdate: "1899-07-21",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		mockRepository.On("CreateTx", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		errService := mockService.Create(mockAppCtx, request)

		// Assert the result
		assert.Nil(t, errService)
	})

	t.Run("CreateAuthor Failed Repository", func(t *testing.T) {
		// Set up input
		request := &entity.UpsertAuthor{
			Name:      "Ernest Hemingway",
			Birthdate: "1899-07-21",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		mockRepository.On("CreateTx", mockAppCtx, mock.Anything, mock.Anything).Return(errors.New("test error"))
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Create(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})

	t.Run("CreateAuthor Validator Failed", func(t *testing.T) {
		// Set up input
		request := &entity.UpsertAuthor{
			Birthdate: "1899-07-21",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Create(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})
}

func TestUpdateAuthor(t *testing.T) {
	mockAppCtx := context.Background()
	t.Run("UpdateAuthor Success", func(t *testing.T) {
		request := &entity.UpsertAuthor{
			Name:      "Ernest Hemingway",
			Birthdate: "1899-07-21",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		mockRepository.On("UpdateTx", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		errService := mockService.Update(mockAppCtx, "123e4567-e89b-12d3-a456-426614174000", request)

		// Assert the result
		assert.Nil(t, errService)
	})

	t.Run("UpdateAuthor Failed Repository", func(t *testing.T) {
		request := &entity.UpsertAuthor{
			Name:      "Ernest Hemingway",
			Birthdate: "1899-07-21",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		mockRepository.On("UpdateTx", mockAppCtx, mock.Anything, mock.Anything).Return(errors.New("test error"))
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Update(mockAppCtx, "123e4567-e89b-12d3-a456-426614174000", request)

		// Assert the result
		assert.NotNil(t, errService)
	})

	t.Run("UpdateAuthor Validator Failed", func(t *testing.T) {
		request := &entity.UpsertAuthor{
			Birthdate: "1899-07-21",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Update(mockAppCtx, "123e4567-e89b-12d3-a456-426614174000", request)

		// Assert the result
		assert.NotNil(t, errService)
	})
}

func TestDeleteAuthor(t *testing.T) {
	mockAppCtx := context.Background()
	t.Run("DeleteAuthor Success", func(t *testing.T) {
		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		mockRepository.On("DeleteByIDTx", mockAppCtx, mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(nil)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		errService := mockService.Delete(mockAppCtx, "123e4567-e89b-12d3-a456-426614174000")

		// Assert the result
		assert.Nil(t, errService)
	})

	t.Run("DeleteAuthor Failed Repository", func(t *testing.T) {
		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		mockRepository.On("DeleteByIDTx", mockAppCtx, mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(errors.New("test error"))
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Delete(mockAppCtx, "123e4567-e89b-12d3-a456-426614174000")

		// Assert the result
		assert.NotNil(t, errService)
	})
}

func TestListAuthor(t *testing.T) {
	mockAppCtx := context.Background()
	req := model.ListReq{
		Page: model.PaginationParam{
			Page:     1,
			PageSize: 1,
		},
		Order: model.OrderParam{
			Order:   "name",
			OrderBy: "asc",
		},
	}
	authors := []*entity.Author{
		{
			Id:        "123e4567-e89b-12d3-a456-426614174000",
			Name:      "Ernest Hemingway",
			Birthdate: nil,
		},
	}

	t.Run("ListAuthor Success", func(t *testing.T) {
		response := &model.PaginationData[entity.Author]{
			Page:             1,
			PageSize:         1,
			TotalPage:        1,
			TotalDataPerPage: 1,
			TotalData:        1,
			Data:             authors,
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		mockRepository.On("FindByPagination", mockAppCtx, mock.Anything, req.Page, req.Order, req.Filter).Return(response, nil)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		result, errService := mockService.List(mockAppCtx, req)

		// Assert the result
		assert.Nil(t, errService)
		assert.NotNil(t, result)
	})

	t.Run("ListAuthor Failed Repository", func(t *testing.T) {
		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		mockRepository.On("FindByPagination", mockAppCtx, mock.Anything, req.Page, req.Order, req.Filter).Return(nil, errors.New("test error"))
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		result, errService := mockService.List(mockAppCtx, req)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})
}

func TestFindOneAuthor(t *testing.T) {
	mockAppCtx := context.Background()
	t.Run("FindOneAuthor Success", func(t *testing.T) {
		response := &entity.Author{
			Id:        "123e4567-e89b-12d3-a456-426614174000",
			Name:      "Ernest Hemingway",
			Birthdate: nil,
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		mockRepository.On("FindByID", mockAppCtx, mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(response, nil)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		result, errService := mockService.FindOne(mockAppCtx, "123e4567-e89b-12d3-a456-426614174000")

		// Assert the result
		assert.Nil(t, errService)
		assert.NotNil(t, result)
	})

	t.Run("FindOneAuthor Failed Repository", func(t *testing.T) {
		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.AuthorRepository)
		mockRepository.On("FindByID", mockAppCtx, mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(nil, errors.New("test error"))
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewAuthorService(gormDB, mockRepository, validate)

		// Call the function under test
		result, errService := mockService.FindOne(mockAppCtx, "123e4567-e89b-12d3-a456-426614174000")

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})
}

package service_test

import (
	"books-api/internal/entity"
	"books-api/internal/mocks"
	"books-api/internal/model"
	service "books-api/internal/services"
	"books-api/pkg/xvalidator"
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func setupSQLMock(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	// Setup SQL mock
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// Setup GORM with the mock DB
	gormDB, gormDBErr := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if gormDBErr != nil {
		t.Fatalf("failed to open GORM connection: %v", gormDBErr)
	}
	return mockSql, gormDB
}

func TestNewBookService(t *testing.T) {
	t.Run("test create BookServiceImpl.NewService BookServiceImpl", func(t *testing.T) {
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		validate, err := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		assert.NotNil(t, mockService)
		assert.Nil(t, err)
	})
}

func TestCreateBook(t *testing.T) {
	mockAppCtx := context.Background()
	t.Run("CreateBook Success", func(t *testing.T) {
		// Set up input
		request := &entity.UpsertBook{
			Title:    "The Great Gatsby",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}

		//mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		mockRepository.On("CreateTx", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		errService := mockService.Create(mockAppCtx, request)

		// Assert the result
		assert.Nil(t, errService)
	})
	t.Run("CreateBook Failed Repository", func(t *testing.T) {
		// Set up input
		request := &entity.UpsertBook{
			Title:    "The Great Gatsby",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}

		//mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		mockRepository.On("CreateTx", mockAppCtx, mock.Anything, mock.Anything).Return(errors.New("test error"))
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Create(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})
	t.Run("CreateBook Validator Failed", func(t *testing.T) {
		// Set up input
		request := &entity.UpsertBook{
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}

		//mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Create(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})
}

func TestListBook(t *testing.T) {
	mockAppCtx := context.Background()
	req := model.ListReq{
		Page: model.PaginationParam{
			Page:     1,
			PageSize: 1,
		},
		Order: model.OrderParam{
			Order:   "title",
			OrderBy: "asc",
		},
	}
	books := []*entity.Book{
		{
			Id:       "0b8d3f3d-d343-4390-964c-4f05c4c803d6",
			Title:    "The Great Gatsby",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		},
	}
	t.Run("ListBook Success", func(t *testing.T) {
		response := &model.PaginationData[entity.Book]{
			Page:             1,
			PageSize:         1,
			TotalPage:        1,
			TotalDataPerPage: 1,
			TotalData:        1,
			Data:             books,
		}
		//mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		mockRepository.On("FindByPagination", mockAppCtx, mock.Anything, req.Page, req.Order, req.Filter).Return(response, nil)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		result, errService := mockService.List(mockAppCtx, req)

		// Assert the result
		assert.Nil(t, errService)
		assert.NotNil(t, result)
	})
	t.Run("ListBook Failed Repository", func(t *testing.T) {
		//mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		mockRepository.On("FindByPagination", mockAppCtx, mock.Anything, req.Page, req.Order, req.Filter).Return(nil, errors.New("test error"))
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		result, errService := mockService.List(mockAppCtx, req)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})
}

func TestFindOneBook(t *testing.T) {
	mockAppCtx := context.Background()
	t.Run("FindOneBook Success", func(t *testing.T) {
		response := &entity.Book{
			Id:       "0b8d3f3d-d343-4390-964c-4f05c4c803d6",
			Title:    "The Great Gatsby",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}
		//mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		mockRepository.On("FindByID", mockAppCtx, mock.Anything, "0b8d3f3d-d343-4390-964c-4f05c4c803d6").Return(response, nil)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		result, errService := mockService.FindOne(mockAppCtx, "0b8d3f3d-d343-4390-964c-4f05c4c803d6")

		// Assert the result
		assert.Nil(t, errService)
		assert.NotNil(t, result)
	})
	t.Run("ListBook Failed Repository", func(t *testing.T) {
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		mockRepository.On("FindByID", mockAppCtx, mock.Anything, "0b8d3f3d-d343-4390-964c-4f05c4c803d6").Return(nil, errors.New("test error"))
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		result, errService := mockService.FindOne(mockAppCtx, "0b8d3f3d-d343-4390-964c-4f05c4c803d6")

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})
}

func TestUpdateBook(t *testing.T) {
	mockAppCtx := context.Background()
	t.Run("UpdateBook Success", func(t *testing.T) {
		request := &entity.UpsertBook{
			Title:    "The Great Gatsby",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}
		//mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		mockRepository.On("UpdateTx", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		errService := mockService.Update(mockAppCtx, "0b8d3f3d-d343-4390-964c-4f05c4c803d6", request)

		// Assert the result
		assert.Nil(t, errService)
	})
	t.Run("UpdateBook Failed Repository", func(t *testing.T) {
		request := &entity.UpsertBook{
			Title:    "The Great Gatsby",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}
		//mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		mockRepository.On("UpdateTx", mockAppCtx, mock.Anything, mock.Anything).Return(errors.New("test error"))
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Update(mockAppCtx, "0b8d3f3d-d343-4390-964c-4f05c4c803d6", request)

		// Assert the result
		assert.NotNil(t, errService)
	})
	t.Run("UpdateBook Validator Failed", func(t *testing.T) {
		// Set up input
		request := &entity.UpsertBook{
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}

		//mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Update(mockAppCtx, "0b8d3f3d-d343-4390-964c-4f05c4c803d6", request)

		// Assert the result
		assert.NotNil(t, errService)
	})
}

func TestDeleteBook(t *testing.T) {
	mockAppCtx := context.Background()
	t.Run("DeleteBook Success", func(t *testing.T) {
		//mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		mockRepository.On("DeleteByIDTx", mockAppCtx, mock.Anything, "0b8d3f3d-d343-4390-964c-4f05c4c803d6").Return(nil)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		errService := mockService.Delete(mockAppCtx, "0b8d3f3d-d343-4390-964c-4f05c4c803d6")

		// Assert the result
		assert.Nil(t, errService)
	})
	t.Run("DeleteBook Failed Repository", func(t *testing.T) {
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.BookRepository)
		mockRepository.On("DeleteByIDTx", mockAppCtx, mock.Anything, "0b8d3f3d-d343-4390-964c-4f05c4c803d6").Return(errors.New("test error"))
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewBookService(gormDB, mockRepository, validate)
		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Delete(mockAppCtx, "0b8d3f3d-d343-4390-964c-4f05c4c803d6")

		// Assert the result
		assert.NotNil(t, errService)
	})
}

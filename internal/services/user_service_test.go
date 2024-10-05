package service_test

import (
	"books-api/internal/entity"
	"books-api/internal/mocks"
	service "books-api/internal/services"
	mocksSignature "books-api/pkg/mocks"
	"books-api/pkg/xvalidator"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	mockAppCtx := context.Background()

	t.Run("RegisterUser Success", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(nil, nil)
		mockRepository.On("CreateTx", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("HashBscryptPassword", request.Password).Return("$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", nil)

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		errService := mockService.Register(mockAppCtx, request)

		// Assert the result
		assert.Nil(t, errService)
	})

	t.Run("RegisterUser Username Exists", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "john_doe",
		}
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(existingUser, nil)

		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Register(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})

	t.Run("RegisterUser Validation Failed", func(t *testing.T) {
		// Set up input (invalid data - missing password)
		request := &entity.UserLogin{
			Username: "john_doe",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Register(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})

	t.Run("RegisterUser HashPassword Failed", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(nil, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("HashBscryptPassword", request.Password).Return("", errors.New("hash error"))

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Register(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})
}

func TestLoginUser(t *testing.T) {
	mockAppCtx := context.Background()

	t.Run("LoginUser Success", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "john_doe",
			Password: "$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", // Hashed password
		}
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(existingUser, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("CheckBscryptPasswordHash", request.Password, existingUser.Password).Return(true)
		mockSignaturer.On("GenerateJWT", existingUser.Username).Return("jwt_token", nil)

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.Nil(t, errService)
		assert.NotNil(t, result)
	})

	t.Run("LoginUser Username Not Found", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(nil, nil)

		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})

	t.Run("LoginUser Invalid Password", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe",
			Password: "WrongPass123!",
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "john_doe",
			Password: "$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", // Hashed password
		}
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(existingUser, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("CheckBscryptPasswordHash", request.Password, existingUser.Password).Return(false)

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})

	t.Run("LoginUser JWT Generation Failed", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "john_doe",
			Password: "$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", // Hashed password
		}
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(existingUser, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("CheckBscryptPasswordHash", request.Password, existingUser.Password).Return(true)
		mockSignaturer.On("GenerateJWT", existingUser.Username).Return("", errors.New("jwt error"))

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})
}

package http

import (
	"books-api/internal/entity"
	"books-api/internal/mocks"
	service "books-api/internal/services"
	"books-api/pkg/exception"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHttpHandler_Register(t *testing.T) {
	t.Run("Register Success", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/register", userHandler.Register)

		// Mock Data
		requestBody := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call
		mockUserService.On("Register", mock.Anything, requestBody).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Register Error - Invalid JSON", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/register", userHandler.Register)

		// Malformed JSON
		malformedJSON := `{"invalid_json"}`
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(malformedJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Register Error - Service Error", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/register", userHandler.Register)

		// Mock Data
		requestBody := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call with error
		mockUserService.On("Register", mock.Anything, requestBody).Return(exception.Internal("error", errors.New("registration failed")))

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHttpHandler_Login(t *testing.T) {
	t.Run("Login Success", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/login", userHandler.Login)

		// Mock Data
		requestBody := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		expectedResponse := &service.UserLoginResponse{
			Username: "john_doe",
			Token:    "jwt_token",
		}

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call
		mockUserService.On("Login", mock.Anything, requestBody).Return(expectedResponse, nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Login Error - Invalid JSON", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/login", userHandler.Login)

		// Malformed JSON
		malformedJSON := `{"invalid_json"}`
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(malformedJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Login Error - Service Error", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/login", userHandler.Login)

		// Mock Data
		requestBody := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call with error
		mockUserService.On("Login", mock.Anything, requestBody).Return(nil, exception.Internal("error", errors.New("login failed")))

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

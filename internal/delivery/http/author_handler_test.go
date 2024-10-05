package http

import (
	"books-api/internal/entity"
	"books-api/internal/mocks"
	services "books-api/internal/services"
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

func TestAuthorHttpHandler_Create(t *testing.T) {
	t.Run("Create Success", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.POST("/authors", authorHandler.Create)

		// Prepare request data
		requestBody := &entity.UpsertAuthor{
			Name:      "F. Scott Fitzgerald",
			Birthdate: "1896-09-24",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/authors", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call
		mockAuthorService.On("Create", mock.Anything, requestBody).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create Error - Invalid JSON", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.POST("/authors", authorHandler.Create)

		// Malformed JSON
		malformedJSON := `{"invalid_json"}`
		req, _ := http.NewRequest("POST", "/authors", bytes.NewBufferString(malformedJSON))
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

	t.Run("Create Error - Service Error", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.POST("/authors", authorHandler.Create)

		// Prepare request data
		requestBody := &entity.UpsertAuthor{
			Name:      "F. Scott Fitzgerald",
			Birthdate: "1896-09-24",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/authors", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call with error
		mockAuthorService.On("Create", mock.Anything, requestBody).Return(exception.Internal("error", errors.New("create failed")))

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAuthorHttpHandler_List(t *testing.T) {
	t.Run("List Success", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.GET("/authors", authorHandler.List)

		// Mock response
		mockResponse := &services.ListAuthorResp{}
		mockAuthorService.On("List", mock.Anything, mock.Anything).Return(mockResponse, nil)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/authors", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("List Error - Service Error", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.GET("/authors", authorHandler.List)

		// Mock error response
		mockAuthorService.On("List", mock.Anything, mock.Anything).Return(nil, exception.Internal("error", errors.New("list failed")))

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/authors", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAuthorHttpHandler_FindOne(t *testing.T) {
	t.Run("FindOne Success", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.GET("/authors/:id", authorHandler.FindOne)

		// Mock response
		authorID := "123e4567-e89b-12d3-a456-426614174000"
		mockResponse := &entity.Author{Id: authorID, Name: "F. Scott Fitzgerald"}
		mockAuthorService.On("FindOne", mock.Anything, authorID).Return(mockResponse, nil)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/authors/"+authorID, nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("FindOne Error", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.GET("/authors/:id", authorHandler.FindOne)

		// Mock error response
		authorID := "invalid-id"
		mockAuthorService.On("FindOne", mock.Anything, authorID).Return(nil, exception.NotFound("author not found"))

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/authors/"+authorID, nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAuthorHttpHandler_Update(t *testing.T) {
	t.Run("Update Success", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.PUT("/authors/:id", authorHandler.Update)

		// Prepare request data
		authorID := "123e4567-e89b-12d3-a456-426614174000"
		requestBody := &entity.UpsertAuthor{
			Name:      "F. Scott Fitzgerald",
			Birthdate: "1896-09-24",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/authors/"+authorID, bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call
		mockAuthorService.On("Update", mock.Anything, authorID, requestBody).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update Error - Invalid JSON", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.PUT("/authors/:id", authorHandler.Update)

		// Malformed JSON
		malformedJSON := `{"invalid_json"}`
		authorID := "123e4567-e89b-12d3-a456-426614174000"
		req, _ := http.NewRequest("PUT", "/authors/"+authorID, bytes.NewBufferString(malformedJSON))
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

	t.Run("Update Error - Service Error", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.PUT("/authors/:id", authorHandler.Update)

		// Prepare request data
		authorID := "123e4567-e89b-12d3-a456-426614174000"
		requestBody := &entity.UpsertAuthor{
			Name:      "F. Scott Fitzgerald",
			Birthdate: "1896-09-24",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/authors/"+authorID, bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call with error
		mockAuthorService.On("Update", mock.Anything, authorID, requestBody).Return(exception.Internal("error", errors.New("update failed")))

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAuthorHttpHandler_Delete(t *testing.T) {
	t.Run("Delete Success", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.DELETE("/authors/:id", authorHandler.Delete)

		// Mock author ID
		authorID := "123e4567-e89b-12d3-a456-426614174000"

		// Mock service call
		mockAuthorService.On("Delete", mock.Anything, authorID).Return(nil)

		// Create HTTP DELETE request
		req, _ := http.NewRequest("DELETE", "/authors/"+authorID, nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Error", func(t *testing.T) {
		r := gin.Default()
		mockAuthorService := new(mocks.AuthorService)
		authorHandler := NewAuthorHTTPHandler(mockAuthorService)

		r.DELETE("/authors/:id", authorHandler.Delete)

		// Mock author ID
		authorID := "123e4567-e89b-12d3-a456-426614174000"

		// Mock service call with error
		mockAuthorService.On("Delete", mock.Anything, authorID).Return(exception.Internal("error", errors.New("delete failed")))

		// Create HTTP DELETE request
		req, _ := http.NewRequest("DELETE", "/authors/"+authorID, nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

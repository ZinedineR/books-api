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

func TestBookHttpHandler_Create(t *testing.T) {
	// Setup router
	t.Run("Positive Case", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		mockUsersService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockUsersService)

		r.POST("/books", bookHandler.Create)
		// Prepare request data
		requestBody := &entity.UpsertBook{
			Title:    "The Great Gatsby",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Create", mock.Anything, requestBody).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

	})
	t.Run("Error service", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		mockUsersService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockUsersService)

		r.POST("/books", bookHandler.Create)
		requestBody := &entity.UpsertBook{
			Title:    "The Great Gatsby",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUsersService.On("Create", mock.Anything, requestBody).Return(exception.Internal("error", errors.New("test error")))

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("Error binding json", func(t *testing.T) {

		t.Parallel()
		r := gin.Default()
		mockUsersService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockUsersService)

		r.POST("/books", bookHandler.Create)
		malformedJSON := `{"very malformed`
		requestBodyBytes, _ := json.Marshal(malformedJSON)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(requestBodyBytes))
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
}

func TestBookHttpHandler_List(t *testing.T) {

	expectResponse := &services.ListBookResp{}

	t.Run("List Success", func(t *testing.T) {

		// Mock
		r := gin.Default()
		mockUsersService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockUsersService)

		// URI
		r.GET("/books", bookHandler.List)

		// Create HTTP Req
		request, _ := http.NewRequest("GET", "/books", nil)
		request.Header.Set("content-type", "application/json")

		// Create Gin Context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = request

		// Mock
		mockUsersService.On("List", mock.Anything, mock.Anything).Return(expectResponse, nil)

		// Perform Request
		r.ServeHTTP(w, request)

		// Check Status Code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("List Failed BookService", func(t *testing.T) {

		// Mock
		r := gin.Default()
		mockUsersService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockUsersService)

		// URI
		r.GET("/books", bookHandler.List)

		// Create HTTP Req
		request, _ := http.NewRequest("GET", "/books", nil)
		request.Header.Set("content-type", "application/json")

		// Create Gin Context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = request

		// Mock
		mockUsersService.On("List", mock.Anything, mock.Anything).Return(nil, exception.Internal("error", errors.New("test error")))
		// Perform Request
		r.ServeHTTP(w, request)

		// Check Status Code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("List Failed ParseParam", func(t *testing.T) {

		// Mock
		r := gin.Default()
		mockUsersService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockUsersService)

		// URI
		r.GET("/books", bookHandler.List)

		// Create HTTP Req
		request, _ := http.NewRequest("GET", "/books?filter=id:try:error", nil)
		request.Header.Set("content-type", "application/json")
		// Create Gin Context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = request

		// Perform Request
		r.ServeHTTP(w, request)

		// Check Status Code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestBookHttpHandler_FindOne(t *testing.T) {
	t.Run("FindOne Success", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockBookService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockBookService)

		r.GET("/books/:id", bookHandler.FindOne)

		// Mock Data
		bookID := "0b8d3f3d-d343-4390-964c-4f05c4c803d6"
		expectedBook := &entity.Book{
			Id:       bookID,
			Title:    "The Great Gatsby",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}

		// Set up the expectation on the mock service
		mockBookService.On("FindOne", mock.Anything, bookID).Return(expectedBook, nil)

		// Create HTTP request
		req, _ := http.NewRequest("GET", "/books/"+bookID, nil)
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
		// Setup
		r := gin.Default()
		mockBookService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockBookService)

		r.GET("/books/:id", bookHandler.FindOne)

		// Mock Data
		bookID := "invalid-id"
		mockBookService.On("FindOne", mock.Anything, bookID).Return(nil, exception.NotFound("book not found"))

		// Create HTTP request
		req, _ := http.NewRequest("GET", "/books/"+bookID, nil)
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

func TestBookHttpHandler_Update(t *testing.T) {
	t.Run("Update Success", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockBookService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockBookService)

		r.PUT("/books/:id", bookHandler.Update)

		// Mock Data
		bookID := "0b8d3f3d-d343-4390-964c-4f05c4c803d6"
		requestBody := &entity.UpsertBook{
			Title:    "Updated Title",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Set up the expectation on the mock service
		mockBookService.On("Update", mock.Anything, bookID, requestBody).Return(nil)

		// Create HTTP request
		req, _ := http.NewRequest("PUT", "/books/"+bookID, bytes.NewBuffer(requestBodyBytes))
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

	t.Run("Update Error - Invalid JSON", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockBookService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockBookService)

		r.PUT("/books/:id", bookHandler.Update)

		// Mock Data
		malformedJSON := `{"invalid_json"}`
		bookID := "0b8d3f3d-d343-4390-964c-4f05c4c803d6"

		// Create HTTP request
		req, _ := http.NewRequest("PUT", "/books/"+bookID, bytes.NewBufferString(malformedJSON))
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
		// Setup
		r := gin.Default()
		mockBookService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockBookService)

		r.PUT("/books/:id", bookHandler.Update)

		// Mock Data
		bookID := "0b8d3f3d-d343-4390-964c-4f05c4c803d6"
		requestBody := &entity.UpsertBook{
			Title:    "Updated Title",
			ISBN:     "978-3-16-148410-0",
			AuthorId: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Set up the expectation on the mock service
		mockBookService.On("Update", mock.Anything, bookID, requestBody).Return(exception.Internal("error", errors.New("update failed")))

		// Create HTTP request
		req, _ := http.NewRequest("PUT", "/books/"+bookID, bytes.NewBuffer(requestBodyBytes))
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

func TestBookHttpHandler_Delete(t *testing.T) {
	t.Run("Delete Success", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockBookService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockBookService)

		r.DELETE("/books/:id", bookHandler.Delete)

		// Mock Data
		bookID := "0b8d3f3d-d343-4390-964c-4f05c4c803d6"

		// Set up the expectation on the mock service
		mockBookService.On("Delete", mock.Anything, bookID).Return(nil)

		// Create HTTP request
		req, _ := http.NewRequest("DELETE", "/books/"+bookID, nil)
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
		// Setup
		r := gin.Default()
		mockBookService := new(mocks.BookService)
		bookHandler := NewBookHTTPHandler(mockBookService)

		r.DELETE("/books/:id", bookHandler.Delete)

		// Mock Data
		bookID := "0b8d3f3d-d343-4390-964c-4f05c4c803d6"

		// Set up the expectation on the mock service
		mockBookService.On("Delete", mock.Anything, bookID).Return(exception.Internal("error", errors.New("delete failed")))

		// Create HTTP request
		req, _ := http.NewRequest("DELETE", "/books/"+bookID, nil)
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

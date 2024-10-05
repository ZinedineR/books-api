package http

import (
	"books-api/pkg/mocks"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignatureHttpHandler_Signature(t *testing.T) {
	t.Run("Signature Success", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockSignatureService := new(mocks.Signaturer)
		signatureHandler := NewSignatureHTTPHandler(mockSignatureService)

		r.POST("/auth/signature", signatureHandler.Signature)

		// Mock data
		httpMethod := "POST"
		bodyJson := `{"key": "value"}`
		token := "auth_token"
		expectedSignature := "generated_signature"

		// Mock SignHMAC512 to return the expected signature
		mockSignatureService.On("SignHMAC512", httpMethod, bodyJson, mock.Anything).Return(expectedSignature, nil)

		// Create HTTP POST request
		req, _ := http.NewRequest(httpMethod, "/auth/signature", bytes.NewBufferString(bodyJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("httpMethod", httpMethod)

		// Create gin context
		w := httptest.NewRecorder()
		w.Header().Set("access_token", token)
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set the access_token manually in the Gin context, as GetToken retrieves it from the context
		//ginCtx.Set("access_token", token)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Signature Error - Invalid HTTP Method", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockSignatureService := new(mocks.Signaturer)
		signatureHandler := NewSignatureHTTPHandler(mockSignatureService)

		r.POST("/auth/signature", signatureHandler.Signature)

		// Create HTTP POST request with invalid HTTP method
		req, _ := http.NewRequest("POST", "/auth/signature", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("httpMethod", "INVALID_METHOD")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Signature Error - Invalid JSON Body", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockSignatureService := new(mocks.Signaturer)
		signatureHandler := NewSignatureHTTPHandler(mockSignatureService)

		r.POST("/auth/signature", signatureHandler.Signature)

		// Mock data
		httpMethod := "POST"
		invalidBodyJson := `{"key": "value"`

		// Create HTTP POST request with invalid JSON body
		req, _ := http.NewRequest(httpMethod, "/auth/signature", bytes.NewBufferString(invalidBodyJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("httpMethod", httpMethod)

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Signature Error - SignHMAC512 Failure", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockSignatureService := new(mocks.Signaturer)
		signatureHandler := NewSignatureHTTPHandler(mockSignatureService)

		r.POST("/auth/signature", signatureHandler.Signature)

		// Mock data
		httpMethod := "POST"
		bodyJson := `{"key": "value"}`
		token := "auth_token"

		// Mock SignHMAC512 to return an error
		mockSignatureService.On("SignHMAC512", httpMethod, bodyJson, mock.Anything).Return("", errors.New("signature generation failed"))

		// Create HTTP POST request
		req, _ := http.NewRequest(httpMethod, "/auth/signature", bytes.NewBufferString(bodyJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("httpMethod", httpMethod)

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set the access_token manually in the Gin context, as GetToken retrieves it from the context
		ginCtx.Set("access_token", token)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

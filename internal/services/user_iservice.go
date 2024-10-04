package service

import (
	"books-api/internal/entity"
	"books-api/pkg/exception"
	"context"
)

type UserService interface {
	// CRUD operations for Book
	Register(
		ctx context.Context, model *entity.UserLogin,
	) *exception.Exception
	Login(ctx context.Context, model *entity.UserLogin) (*UserLoginResponse, *exception.Exception)
}

type UserLoginResponse struct {
	Username string `json:"username" example:"john_doe"`
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"` // JWT token example
}

package route

import (
	"books-api/internal/delivery/http"
	api "books-api/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	App              *gin.Engine
	BookHandler      *http.BookHTTPHandler
	AuthorHandler    *http.AuthorHTTPHandler
	UserHandler      *http.UserHTTPHandler
	SignatureHandler *http.SignatureHTTPHandler
	AuthMiddleware   *api.AuthMiddleware
}

func (h *Router) Setup() {
	h.App.Use(h.AuthMiddleware.ErrorHandler)
	guestApi := h.App.Group("/auth")
	{
		guestApi.POST("/register", h.UserHandler.Register)
		guestApi.POST("/login", h.UserHandler.Login)
	}
	signatureAPI := h.App.Group("/auth/signature")
	signatureAPI.Use(h.AuthMiddleware.JWTAuthentication)
	{
		signatureAPI.POST("", h.SignatureHandler.Signature)
	}
	userApi := h.App.Group("")
	userApi.Use(h.AuthMiddleware.JWTAuthentication, h.AuthMiddleware.SignatureAuthentication)
	{
		booksApi := userApi.Group("/books")
		{
			booksApi.POST("", h.BookHandler.Create)
			booksApi.GET("", h.BookHandler.List)
			booksApi.GET("/:id", h.BookHandler.FindOne)
			booksApi.PUT("/:id", h.BookHandler.Update)
			booksApi.DELETE("/:id", h.BookHandler.Delete)
		}
		authorApi := userApi.Group("/authors")
		{
			authorApi.POST("", h.AuthorHandler.Create)
			authorApi.GET("", h.AuthorHandler.List)
			authorApi.GET("/:id", h.AuthorHandler.FindOne)
			authorApi.PUT("/:id", h.AuthorHandler.Update)
			authorApi.DELETE("/:id", h.AuthorHandler.Delete)
		}
	}
}

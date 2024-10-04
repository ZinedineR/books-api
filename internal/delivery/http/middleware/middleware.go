package api

import "books-api/internal/delivery/http"

type Middleware struct {
	http.Handler
}

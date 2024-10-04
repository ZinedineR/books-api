package api

import "boiler-plate-clean/internal/delivery/http"

type Middleware struct {
	http.Handler
}

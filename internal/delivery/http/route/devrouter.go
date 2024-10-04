package route

import (
	_ "books-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (h *Router) setupDevRouter() {

}

func (h *Router) SwaggerRouter() {
	h.App.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

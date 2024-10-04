package route

import (
	"boiler-plate-clean/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

type Router struct {
	App            *gin.Engine
	ExampleHandler *http.ExampleHTTPHandler
}

func (h *Router) Setup() {
	api := h.App.Group("")
	{

		//Example Routes
		campaignApi := api.Group("/campaign")
		//campaignApi.Use(h.RequestMiddleware.RequestHeader)
		{
			campaignApi.POST("/", h.ExampleHandler.Create)
			campaignApi.GET("/select", h.ExampleHandler.Find)
			campaignApi.GET("/:id", h.ExampleHandler.FindOne)
			campaignApi.PUT("/:id", h.ExampleHandler.Update)
			campaignApi.DELETE("/:id", h.ExampleHandler.Delete)
		}
	}
}

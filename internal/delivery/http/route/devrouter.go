package route

import (
	"github.com/rakyll/statik/fs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	//_ "boiler-plate-clean/statik"
)

func (h *Router) setupDevRouter() {

}

func (h *Router) SwaggerRouter() {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	//h.router.Use(static.Serve("/"+h.base.AppConfig.AppConfig.AppName+"/api/"+h.base.AppConfig.AppConfig.AppVersion+"/static", static.LocalFile("./docs", true)))
	//h.router.Static("/"+h.base.AppConfig.AppConfig.AppName+"/api/"+h.base.AppConfig.AppConfig.AppVersion+"/static/*any", "./docs")
	h.App.StaticFS("/static", statikFS)
	h.App.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/static/openapi.yaml")))
}

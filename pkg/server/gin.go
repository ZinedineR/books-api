package server

import (
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type GinConfig struct {
	HttpPort     string
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

type GinServer struct {
	App  *gin.Engine
	Port string
}

func NewGinServer(conf *GinConfig) *GinServer {
	app := gin.New()

	app.Use(gin.Recovery())
	app.Use(sloggin.New(slog.Default()))
	app.Use(gin.Logger())
	app.Use(cors.New(cors.Config{
		AllowOrigins: conf.AllowOrigins,
		AllowMethods: conf.AllowMethods,
		AllowHeaders: conf.AllowHeaders,
	}))

	return &GinServer{
		App:  app,
		Port: conf.HttpPort,
	}
}

func (s *GinServer) Start() error {
	return s.App.Run(":" + s.Port)
}

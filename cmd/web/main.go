package main

import (
	"books-api/config"
	"books-api/internal/delivery/http"
	api "books-api/internal/delivery/http/middleware"
	"books-api/internal/delivery/http/route"
	"books-api/internal/repository"
	services "books-api/internal/services"
	"books-api/migration"
	"books-api/pkg/database"
	"books-api/pkg/httpclient"
	"books-api/pkg/logger"
	"books-api/pkg/server"
	"books-api/pkg/signature"
	"books-api/pkg/xvalidator"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	sqlClientRepo *database.Database
)

// @title           Books-API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9004
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	validate, _ := xvalidator.NewValidator()
	conf := config.InitAppConfig(validate)
	logger.SetupLogger(&logger.Config{
		AppENV:  conf.AppEnvConfig.AppEnv,
		LogPath: conf.AppEnvConfig.LogFilePath,
		Debug:   conf.AppEnvConfig.AppDebug,
	})
	initInfrastructure(conf)
	ginServer := server.NewGinServer(&server.GinConfig{
		HttpPort:     conf.AppEnvConfig.HttpPort,
		AllowOrigins: conf.AppEnvConfig.AllowOrigins,
		AllowMethods: conf.AppEnvConfig.AllowMethods,
		AllowHeaders: conf.AppEnvConfig.AllowHeaders,
	})
	// external
	signaturer := signature.NewSignature(conf.AuthConfig.JwtSecretAccessToken)
	// repository
	booksRepository := repository.NewBookSQLRepository()
	authorRepository := repository.NewAuthorSQLRepository()
	userRepository := repository.NewUserSQLRepository()

	// service
	booksService := services.NewBookService(sqlClientRepo.GetDB(), booksRepository, validate)
	authorService := services.NewAuthorService(sqlClientRepo.GetDB(), authorRepository, validate)
	userService := services.NewUserService(sqlClientRepo.GetDB(), userRepository, signaturer, validate)
	// Handler
	authMiddleware := api.NewAuthMiddleware(signaturer)
	booksHandler := http.NewBookHTTPHandler(booksService)
	authorHandler := http.NewAuthorHTTPHandler(authorService)
	userHandler := http.NewUserHTTPHandler(userService)

	router := route.Router{
		App:            ginServer.App,
		BookHandler:    booksHandler,
		AuthorHandler:  authorHandler,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
	router.Setup()
	router.SwaggerRouter()
	echan := make(chan error)
	go func() {
		echan <- ginServer.Start()
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		slog.Info("signal terminated detected")
	case err := <-echan:
		slog.Error("Failed to start http server", err)
	}
}

func initInfrastructure(config *config.Config) {
	//initPostgreSQL()
	sqlClientRepo = initSQL(config)
}

func initSQL(conf *config.Config) *database.Database {
	db := database.NewDatabase(conf.DatabaseConfig.Dbservice, &database.Config{
		DbHost:   conf.DatabaseConfig.Dbhost,
		DbUser:   conf.DatabaseConfig.Dbuser,
		DbPass:   conf.DatabaseConfig.Dbpassword,
		DbName:   conf.DatabaseConfig.Dbname,
		DbPort:   strconv.Itoa(conf.DatabaseConfig.Dbport),
		DbPrefix: conf.DatabaseConfig.DbPrefix,
	})
	if conf.IsStaging() {
		migration.AutoMigration(db)
	}
	return db
}

func initHttpclient() httpclient.Client {
	httpClientFactory := httpclient.New()
	httpClient := httpClientFactory.CreateClient()
	return httpClient
}

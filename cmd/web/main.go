package main

import (
	"boiler-plate-clean/config"
	"boiler-plate-clean/internal/delivery/http"
	"boiler-plate-clean/internal/delivery/http/route"
	"boiler-plate-clean/internal/gateway/messaging"
	"boiler-plate-clean/internal/repository"
	services "boiler-plate-clean/internal/services"
	"boiler-plate-clean/migration"
	"boiler-plate-clean/pkg/server"
	kafkaserver "github.com/RumbiaID/pkg-library/app/pkg/broker/kafkaservice"
	"github.com/RumbiaID/pkg-library/app/pkg/database"
	"github.com/RumbiaID/pkg-library/app/pkg/httpclient"
	"github.com/RumbiaID/pkg-library/app/pkg/logger"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	httpClient      httpclient.Client
	sqlClientRepo   *database.Database
	kafkaDialer     *kafkaserver.KafkaService
	exampleProducer messaging.ExampleProducer
)

// @title           Pigeon
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/notificationsvc/api/v1

// @securityDefinitions.basic  BasicAuth

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

	// repository
	exampleRepository := repository.NewExampleSQLRepository()

	// external api
	//gotifySvcExternalAPI := externalapi.NewExampleExternalImpl(conf, httpClient)

	// producer

	exampleProducer = messaging.NewExampleKafkaProducerImpl(kafkaDialer, conf.KafkaConfig.KafkaTopicEmail)

	// service
	exampleService := services.NewExampleService(sqlClientRepo.GetDB(), exampleRepository, validate)
	// Handler
	exampleHandler := http.NewExampleHTTPHandler(exampleService)

	router := route.Router{
		App:            ginServer.App,
		ExampleHandler: exampleHandler,
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

	kafkaDialer = initKafka(config)

	sqlClientRepo = initSQL(config)

	httpClient = initHttpclient()
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
	if conf.UseReplica() {
		db.CqrsDB(conf.DatabaseConfig.Dbservice, &database.Config{
			DbHost: conf.DatabaseReplicaConfig.Dbreplicahost,
			DbUser: conf.DatabaseReplicaConfig.Dbreplicauser,
			DbPass: conf.DatabaseReplicaConfig.Dbreplicapassword,
			DbName: conf.DatabaseReplicaConfig.Dbreplicaname,
			DbPort: strconv.Itoa(conf.DatabaseReplicaConfig.Dbreplicaport),
		})
	}
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

func initKafka(config *config.Config) *kafkaserver.KafkaService {
	kafkaDialer := kafkaserver.New(&kafkaserver.Config{
		SecurityProtocol: config.KafkaConfig.KafkaSecurityProtocol,
		Brokers:          config.KafkaConfig.KafkaBroker,
		Username:         config.KafkaConfig.KafkaUsername,
		Password:         config.KafkaConfig.KafkaPassword,
	})
	return kafkaDialer
}

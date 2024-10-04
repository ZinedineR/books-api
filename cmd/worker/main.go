package main

import (
	"boiler-plate-clean/config"
	"boiler-plate-clean/internal/delivery/messaging"
	"context"
	kafkaserver "github.com/RumbiaID/pkg-library/app/pkg/broker/kafkaservice"
	"github.com/RumbiaID/pkg-library/app/pkg/logger"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	kafkaService *kafkaserver.KafkaService
)

func main() {
	validate, _ := xvalidator.NewValidator()
	conf := config.InitConsumerConfig(validate)
	logger.SetupLogger(&logger.Config{
		AppENV:  conf.AppEnvConfig.AppEnv,
		LogPath: conf.AppEnvConfig.LogFilePath,
		Debug:   conf.AppEnvConfig.AppDebug,
	})

	ctx, cancel := context.WithCancel(context.Background())
	//ctx, cancel := context.WithCancel(context.Background())
	// repository

	// external api
	//httpClientFactory := httpclient.New()
	//httpClient := httpClientFactory.CreateClient()

	//Handler
	exampleHandler := messaging.NewExampleConsumer()

	kafkaService = kafkaserver.New(&kafkaserver.Config{
		SecurityProtocol: conf.KafkaConfig.KafkaSecurityProtocol,
		Brokers:          conf.KafkaConfig.KafkaBroker,
		Username:         conf.KafkaConfig.KafkaUsername,
		Password:         conf.KafkaConfig.KafkaPassword,
	})
	go messaging.ConsumeKafkaTopic(ctx, kafkaService, conf.KafkaConfig.KafkaTopicNotification, conf.KafkaConfig.KafkaGroupId, exampleHandler.ConsumeKafka)

	slog.Info("Worker is running")

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	stop := false
	for !stop {
		select {
		case s := <-terminateSignals:
			slog.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
			cancel()
			stop = true
		}
	}

	time.Sleep(5 * time.Second) // wait for all consumers to finish processing
}

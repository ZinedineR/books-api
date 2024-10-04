package kafkaservice

import (
	"crypto/tls"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"log/slog"
	"os"
)

type Config struct {
	SecurityProtocol string
	Brokers          []string
	Username         string
	Password         string
}

type KafkaService struct {
	brokers     []string
	mechanism   sasl.Mechanism
	tls         *tls.Config
	errorLogger kafka.Logger
	logger      kafka.Logger
}

func New(config *Config) *KafkaService {
	var mechanism sasl.Mechanism
	var tlsConfig *tls.Config

	if config.SecurityProtocol == "SCRAM_SHA_256" {
		var err error
		mechanism, err = scram.Mechanism(scram.SHA256, config.Username, config.Password)
		if err != nil {
			slog.Error("can't setup kafka server", "error", err)
			os.Exit(1)
		}
		tlsConfig = &tls.Config{}
	} else if config.SecurityProtocol == "SCRAM_SHA_512" {
		var err error
		mechanism, err = scram.Mechanism(scram.SHA512, config.Username, config.Password)
		if err != nil {
			slog.Error("can't setup kafka server", "error", err)
			os.Exit(1)
		}
		tlsConfig = &tls.Config{}
	} else if config.SecurityProtocol == "SASL_SSL" {
		mechanism = plain.Mechanism{
			Username: config.Username,
			Password: config.Password,
		}
		tlsConfig = &tls.Config{}
	} else if config.SecurityProtocol == "SASL_PLAIN" {
		mechanism = plain.Mechanism{
			Username: config.Username,
			Password: config.Password,
		}
		tlsConfig = nil
	} else if config.SecurityProtocol == "PLAIN" {
		mechanism = nil
		tlsConfig = nil
	} else {
		slog.Error("can't setup kafka server", "error", "invalid security protocol")
		os.Exit(1)
	}

	slog.Debug(fmt.Sprintf("kafka security protocol %s", config.SecurityProtocol))
	slog.Debug("kafka module initialized")

	return &KafkaService{
		brokers:   config.Brokers,
		mechanism: mechanism,
		tls:       tlsConfig,
		errorLogger: kafka.LoggerFunc(func(message string, args ...interface{}) {
			slog.Error(fmt.Sprintf(message, args...))
		}),
		// logger: kafka.LoggerFunc(func(message string, args ...interface{}) {
		// 	slog.Debug(fmt.Sprintf(message, args...))
		// }),
	}
}

func (k *KafkaService) NewReader(topic string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.brokers,
		Topic:   topic,
		GroupID: groupID,
		Dialer: &kafka.Dialer{
			SASLMechanism: k.mechanism,
			TLS:           k.tls,
		},
		ErrorLogger: k.errorLogger,
		Logger:      k.logger,
	})
}

func (k *KafkaService) NewWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(k.brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			SASL: k.mechanism,
			TLS:  k.tls,
		},
		ErrorLogger: k.errorLogger,
		Logger:      k.logger,
	}
}

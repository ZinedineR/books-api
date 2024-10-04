package messaging

import (
	"context"
	kafkaserver "github.com/RumbiaID/pkg-library/app/pkg/broker/kafkaservice"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type ConsumerKafkaHandler func(ctx context.Context, message *kafka.Message) error

func ConsumeKafkaTopic(
	ctx context.Context, kafkaService *kafkaserver.KafkaService, topic, groupid string,
	handler ConsumerKafkaHandler,
) {
	reader := kafkaService.NewReader(topic, groupid)

	run := true

	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			message, err := reader.ReadMessage(ctx)
			if err == nil {
				err := handler(ctx, &message)
				if err != nil {
					slog.Error("Failed to process message: ", slog.Any("error", err))
				} else {
					err = reader.CommitMessages(ctx, message)
					if err != nil {
						slog.Error("Failed to commit message: ", slog.Any("error", err))
					}
				}
			} else {
				slog.Warn("Consumer error: ", slog.Any("error", err))
			}
		}
	}

	slog.Info("Closing consumer for topic : ", slog.String("topic", topic))
	err := reader.Close()
	if err != nil {
		panic(err)
	}
}

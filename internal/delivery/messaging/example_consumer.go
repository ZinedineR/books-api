package messaging

import (
	"boiler-plate-clean/internal/model"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type ExampleConsumer struct {
}

func NewExampleConsumer() *ExampleConsumer {
	return &ExampleConsumer{}
}

func (c ExampleConsumer) ConsumeKafka(ctx context.Context, message *kafka.Message) error {
	exampleEvent := new(model.ExampleMessage)
	if err := json.Unmarshal(message.Value, exampleEvent); err != nil {
		slog.Error("error unmarshalling example event", slog.String("error", err.Error()))
		return err
	}
	slog.Info("Received topic example with event", slog.Any("example", exampleEvent))
	return nil
}

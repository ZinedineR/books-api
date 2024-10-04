package messaging

import (
	"books-api/internal/model"
	"context"
)

type ExampleProducer interface {
	GetTopic() string
	Send(ctx context.Context, order ...*model.ExampleMessage) error
}

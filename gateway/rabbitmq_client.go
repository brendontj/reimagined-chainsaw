package gateway

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient interface {
	Connect() error
	PublishMessage(ctx context.Context, msg string) error
	ConsumeMessages() (<-chan amqp.Delivery, error)
	Close()
}

package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	URL                string
	PublisherQueueName string
	ConsumerQueueName  string
	conn               *amqp.Connection
	channel            *amqp.Channel
	publisherQueue     amqp.Queue
	consumerQueue      amqp.Queue
}

func (c *Client) Connect() error {
	//Todo: Change it to put this value in a env file
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	c.conn = conn
	if err = c.openChannels(); err != nil {
		return err
	}
	return nil
}

func (c *Client) PublishMessage(ctx context.Context, msg string) error {
	return c.channel.PublishWithContext(ctx,
		"",
		c.publisherQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
}

func (c *Client) ConsumeMessages() (<-chan amqp.Delivery, error) {
	messages, err := c.channel.Consume(
		c.consumerQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *Client) Close() {
	defer func() {
		_ = c.channel.Close()
		_ = c.conn.Close()
	}()
}

func (c *Client) openChannels() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	c.channel = ch

	q, err := c.channel.QueueDeclare(
		c.PublisherQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	c.publisherQueue = q

	cq, err := c.channel.QueueDeclare(
		c.ConsumerQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	c.consumerQueue = cq
	return nil
}

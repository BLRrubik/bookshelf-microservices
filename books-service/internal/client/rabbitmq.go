package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const QueueImageCompress = "image_compress"

type ImageCompressMessage struct {
	BookID       string `json:"book_id"`
	CoverID      string `json:"cover_id"`
	OriginalPath string `json:"original_path"`
}

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *RabbitMQClient) DeclareQueue(name string) error {
	q, err := c.channel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return err
	}

	c.queue = q.Name

	return nil
}

func (c *RabbitMQClient) Close() error {
	if err := c.channel.Close(); err != nil {
		_ = c.conn.Close()

		return err
	}

	return c.conn.Close()
}

func (c *RabbitMQClient) HealthCheck() error {
	if c.channel.IsClosed() {
		return errors.New("channel is closed")
	}

	if c.conn.IsClosed() {
		return errors.New("connection is closed")
	}

	return nil
}

func (c *RabbitMQClient) PublishImageCompress(ctx context.Context, msg ImageCompressMessage) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal image compress message: %w", err)
	}

	if err = c.publish(ctx, c.queue, bytes); err != nil {
		return fmt.Errorf("publish image compress message: %w", err)
	}

	return nil
}

func (c *RabbitMQClient) publish(ctx context.Context, queue string, body []byte) error {
	return c.channel.PublishWithContext(
		ctx,
		amqp.DefaultExchange,
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

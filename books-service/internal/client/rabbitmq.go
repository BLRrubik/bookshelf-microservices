package client

import (
	"context"
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
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
	logger  *zap.Logger
}

func NewRabbitMQClient(url string, logger *zap.Logger) (*RabbitMQClient, error) {
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
		logger:  logger,
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
		c.logger.Error("failed to marshal image compress message", zap.String("book_id", msg.BookID), zap.Error(err))

		return err
	}

	if err = c.publish(ctx, c.queue, bytes); err != nil {
		c.logger.Error("failed to publish image compress message", zap.String("book_id", msg.BookID), zap.Error(err))

		return err
	}

	c.logger.Info("published image compress message", zap.String("book_id", msg.BookID), zap.String("cover_id", msg.CoverID))

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

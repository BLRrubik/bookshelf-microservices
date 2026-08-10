package client

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

const QueueImageCompress = "image_compress"

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
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
	_, err := c.channel.QueueDeclare(name, true, false, false, false, nil)

	return err
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

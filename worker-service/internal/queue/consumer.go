package queue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrTemporary = errors.New("temporary error")

type HandlerFunc func(body []byte) error

type Consumer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	handlers map[string]HandlerFunc
}

func NewConsumer(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:     conn,
		channel:  ch,
		handlers: make(map[string]HandlerFunc),
	}, nil
}

func (c *Consumer) RegisterHandler(queue string, handler HandlerFunc) {
	c.handlers[queue] = handler
}

func (c *Consumer) Start() error {
	for queue, handler := range c.handlers {
		if err := c.channel.Qos(1, 0, false); err != nil {
			return err
		}

		deliveries, err := c.channel.ConsumeWithContext(
			context.Background(),
			queue,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		go c.consume(queue, handler, deliveries)
	}

	return nil
}

func (c *Consumer) consume(queue string, handler HandlerFunc, msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		if err := handler(msg.Body); err != nil {
			var requeue bool
			if errors.Is(err, ErrTemporary) {
				requeue = true
			}

			if err = msg.Nack(false, requeue); err != nil {
				slog.Error("Failed to nack message from queue %s: %s", queue, err)
			}

			continue
		}

		if err := msg.Ack(false); err != nil {
			slog.Error(fmt.Sprintf("Failed to ack message from queue %s: %s", queue, err))
		}
	}
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		_ = c.conn.Close()

		return err
	}

	return c.conn.Close()
}

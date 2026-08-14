package queue

import (
	"context"
	"errors"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var ErrTemporary = errors.New("temporary error")

type HandlerFunc func(body []byte) error

type Consumer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	handlers map[string]HandlerFunc
	wg       sync.WaitGroup
	logger   *zap.Logger
}

func NewConsumer(url string, logger *zap.Logger) (*Consumer, error) {
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
		logger:   logger,
	}, nil
}

func (c *Consumer) RegisterHandler(queue string, handler HandlerFunc) {
	c.handlers[queue] = handler
}

func (c *Consumer) Start(ctx context.Context) error {
	for queue, handler := range c.handlers {
		if err := c.channel.Qos(1, 0, false); err != nil {
			return err
		}

		deliveries, err := c.channel.ConsumeWithContext(
			ctx,
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

		c.wg.Add(1)
		go c.consume(ctx, queue, handler, deliveries)
	}

	return nil
}

func (c *Consumer) consume(ctx context.Context, queue string, handler HandlerFunc, msgs <-chan amqp.Delivery) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}

			if err := c.processMessage(queue, handler, msg); err != nil {
				c.logger.Error("error processing message", zap.String("queue", queue), zap.Error(err))
			}
		}
	}
}

func (c *Consumer) Wait(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Consumer) processMessage(queue string, handler HandlerFunc, msg amqp.Delivery) error {
	if err := handler(msg.Body); err != nil {
		var requeue bool
		if errors.Is(err, ErrTemporary) {
			requeue = true
		}

		if nackErr := msg.Nack(false, requeue); nackErr != nil {
			c.logger.Error("failed to nack message", zap.String("queue", queue), zap.Error(nackErr))
		}

		return err
	}

	if err := msg.Ack(false); err != nil {
		c.logger.Error("failed to ack message", zap.String("queue", queue), zap.Error(err))
	}

	return nil
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		_ = c.conn.Close()

		return err
	}

	return c.conn.Close()
}

func (c *Consumer) HealthCheck() error {
	if c.channel.IsClosed() {
		return errors.New("channel is closed")
	}

	if c.conn.IsClosed() {
		return errors.New("connection is closed")
	}

	return nil
}

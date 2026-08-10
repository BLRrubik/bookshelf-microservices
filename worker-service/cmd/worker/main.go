package main

import (
	"bookshelf/worker-service/internal/config"
	"bookshelf/worker-service/internal/queue"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	slog.Info("worker-service starting...")

	slog.Info("Loading configuration...")
	cfg := config.Load()

	consumer, err := queue.NewConsumer(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	consumer.RegisterHandler("test_queue", func(body []byte) error {
		fmt.Println(string(body))

		return nil
	})

	consumer.Start()

	term := make(chan os.Signal, 1)

	signal.Notify(term, os.Interrupt)

	<-term
}

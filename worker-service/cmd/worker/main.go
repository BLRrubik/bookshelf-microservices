package main

import (
	"bookshelf/worker-service/internal/config"
	"bookshelf/worker-service/internal/handler"
	"bookshelf/worker-service/internal/queue"
	"bookshelf/worker-service/internal/repository"
	"bookshelf/worker-service/internal/storage"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	minioStorage, err := storage.NewMinIOStorage(
		cfg.MinIOURL,
		cfg.MinIOAccess,
		cfg.MinIOSecret,
		cfg.MinioBucket,
		cfg.MinioPublicEndpoint,
		cfg.MinioUseSSL,
	)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewCoverRepository(db)

	imageHandler := handler.NewImageHandler(minioStorage, repo)

	consumer.RegisterHandler("test_queue", func(body []byte) error {
		return imageHandler.HandleImageCompress(body)
	})

	consumer.Start()

	term := make(chan os.Signal, 1)

	signal.Notify(term, os.Interrupt)

	<-term
}

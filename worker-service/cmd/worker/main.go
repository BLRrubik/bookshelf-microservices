package main

import (
	"bookshelf/worker-service/internal/api"
	"bookshelf/worker-service/internal/config"
	"bookshelf/worker-service/internal/handler"
	"bookshelf/worker-service/internal/queue"
	"bookshelf/worker-service/internal/repository"
	"bookshelf/worker-service/internal/storage"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	slog.Info("worker-service starting...")

	ctx, cancel := context.WithCancel(context.Background())

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

	consumer.RegisterHandler("image_compress", func(body []byte) error {
		return imageHandler.HandleImageCompress(body)
	})

	if err = consumer.Start(ctx); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5175", "http://localhost:3003"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	httpHandler := api.NewHealthHandler(db, consumer, minioStorage)

	r.Get("/health", httpHandler.Health)
	r.Get("/ready", httpHandler.Ready)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         cfg.Port,
		Handler:      r,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	term := make(chan os.Signal, 1)

	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	<-term

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err = consumer.Wait(shutdownCtx); err != nil {
		log.Printf("consumer shutdown wait error: %v", err)
	}

	consumer.Close()

	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

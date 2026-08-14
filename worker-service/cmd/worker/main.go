package main

import (
	"bookshelf/worker-service/internal/api"
	"bookshelf/worker-service/internal/config"
	"bookshelf/worker-service/internal/handler"
	"bookshelf/worker-service/internal/logger"
	"bookshelf/worker-service/internal/queue"
	"bookshelf/worker-service/internal/repository"
	"bookshelf/worker-service/internal/storage"
	"context"
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
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel)
	defer log.Sync()

	log.Info("config loaded", zap.String("port", cfg.Port))

	ctx, cancel := context.WithCancel(context.Background())

	consumer, err := queue.NewConsumer(cfg.RabbitMQURL, log.Named("consumer"))
	if err != nil {
		log.Fatal("failed to connect to rabbitmq", zap.Error(err))
	}
	defer consumer.Close()

	log.Info("connected to rabbitmq")

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	log.Info("connected to database")

	minioStorage, err := storage.NewMinIOStorage(
		cfg.MinIOURL,
		cfg.MinIOAccess,
		cfg.MinIOSecret,
		cfg.MinioBucket,
		cfg.MinioPublicEndpoint,
		cfg.MinioUseSSL,
		log.Named("minio"),
	)
	if err != nil {
		log.Fatal("failed to create minio client", zap.Error(err))
	}

	log.Info("connected to minio")

	repo := repository.NewCoverRepository(db, log.Named("cover_repository"))

	imageHandler := handler.NewImageHandler(minioStorage, repo, log.Named("image_handler"))

	consumer.RegisterHandler("image_compress", func(body []byte) error {
		return imageHandler.HandleImageCompress(body)
	})

	if err = consumer.Start(ctx); err != nil {
		log.Fatal("failed to start consumer", zap.Error(err))
	}

	log.Info("consumer started")

	r := chi.NewRouter()

	r.Use(api.RequestLogger(log.Named("http")))
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
		log.Info("server listening", zap.String("addr", cfg.Port))

		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	term := make(chan os.Signal, 1)

	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	<-term

	log.Info("shutdown signal received")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err = consumer.Wait(shutdownCtx); err != nil {
		log.Error("consumer shutdown wait error", zap.Error(err))
	}

	consumer.Close()

	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", zap.Error(err))
	} else {
		log.Info("graceful shutdown complete")
	}
}

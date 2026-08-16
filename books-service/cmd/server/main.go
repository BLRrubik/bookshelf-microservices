package main

import (
	"bookshelf/books-service/internal/client"
	"bookshelf/books-service/internal/config"
	"bookshelf/books-service/internal/handler"
	"bookshelf/books-service/internal/logger"
	"bookshelf/books-service/internal/repository"
	"bookshelf/books-service/internal/service"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	log := logger.New(cfg.LogLevel)
	defer log.Sync()

	log.Info("config loaded", zap.String("port", cfg.Port))

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	log.Info("connected to database")

	rabbitMQClient, err := client.NewRabbitMQClient("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatal("failed to connect to rabbitmq", zap.Error(err))
	}
	defer rabbitMQClient.Close()

	log.Info("connected to rabbitmq")

	if err = rabbitMQClient.DeclareQueue(client.QueueImageCompress); err != nil {
		log.Fatal("failed to declare queue", zap.Error(err))
	}

	minioClient, err := client.NewMinIOClient(
		cfg.MinIOEndpoint,
		cfg.MinIOAccessKey,
		cfg.MinIOSecretKey,
		cfg.MinioBucket,
		cfg.MinIOPublicEndpoint,
		cfg.MinIOUseSSL,
	)
	if err != nil {
		log.Fatal("failed to create minio client", zap.Error(err))
	}

	if err = minioClient.EnsureBucket(ctx); err != nil {
		log.Fatal("failed to ensure bucket", zap.Error(err))
	}

	log.Info("connected to minio")

	_ = minioClient

	authClient := client.NewAuthClient(
		"http://localhost:8081",
		time.Second*5,
		5,
		time.Millisecond*100,
		cfg.ServiceKey,
		log.Named("auth_client"),
	)

	repos := repository.New(db)
	services := service.New(repos, authClient, minioClient, rabbitMQClient, log.Named("service"))
	handlers := handler.NewHandler(services, db, authClient, rabbitMQClient, minioClient, log.Named("handler"))

	r := chi.NewRouter()

	r.Use(handler.RequestLogger(log.Named("http")))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	handlers.RegisterRoutes(r)

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

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	} else {
		log.Info("graceful shutdown complete")
	}
}

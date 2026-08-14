package main

import (
	"bookshelf/auth-service/internal/config"
	"bookshelf/auth-service/internal/handler"
	"bookshelf/auth-service/internal/logger"
	"bookshelf/auth-service/internal/repository"
	"bookshelf/auth-service/internal/service"
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

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	log.Info("connected to database")

	userRepo := repository.NewUserRepository(db, log.Named("user_repository"))
	userService := service.NewUserService(userRepo, cfg.JWTSecret, log.Named("user_service"))
	handlers := handler.NewHandler(userService, db, cfg.JWTSecret, log.Named("handler"))

	r := chi.NewRouter()

	r.Use(handler.RequestLogger(log.Named("http")))
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

	handlers.RegisterRoutes(r, cfg.ServiceKey)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	} else {
		log.Info("graceful shutdown complete")
	}
}

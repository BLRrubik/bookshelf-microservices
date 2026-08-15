package main

import (
	"bookshelf/api-gateway/internal/config"
	"bookshelf/api-gateway/internal/handler"
	"bookshelf/api-gateway/internal/logger"
	"bookshelf/api-gateway/internal/middleware"
	"bookshelf/api-gateway/internal/proxy"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	log := logger.New(zap.DebugLevel.String())
	defer log.Sync()

	proxyService := proxy.New(cfg.AuthServiceURL, cfg.BooksServiceURL)
	handlers := handler.NewHandler(proxyService)

	r := chi.NewRouter()

	r.Use(middleware.RequestLogger(log.Named("http")))
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.RequestID())
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

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

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	} else {
		log.Info("graceful shutdown complete")
	}
}

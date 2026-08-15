package main

import (
	"bookshelf/api-gateway/internal/config"
	"bookshelf/api-gateway/internal/handler"
	"bookshelf/api-gateway/internal/middleware"
	"bookshelf/api-gateway/internal/proxy"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	var logHandler slog.Handler
	if cfg.LogFormat == "text" {
		logHandler = slog.NewTextHandler(os.Stdout, nil)
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	}
	slog.SetDefault(slog.New(logHandler))

	proxyService := proxy.New(cfg.AuthServiceURL, cfg.BooksServiceURL)
	handlers := handler.NewHandler(proxyService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logging)
	r.Use(middleware.Cors)

	handlers.RegisterRoutes(r)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         cfg.Port,
		Handler:      r,
	}

	go func() {
		slog.Info("server listening", "addr", cfg.Port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	} else {
		slog.Info("graceful shutdown complete")
	}
}

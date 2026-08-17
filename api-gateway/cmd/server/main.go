package main

import (
	"bookshelf/api-gateway/internal/cache"
	"bookshelf/api-gateway/internal/config"
	"bookshelf/api-gateway/internal/handler"
	"bookshelf/api-gateway/internal/logger"
	"bookshelf/api-gateway/internal/middleware"
	"bookshelf/api-gateway/internal/proxy"
	"bookshelf/api-gateway/internal/ratelimit"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	lgr := logger.New("info")

	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalln(err)
	}

	redisClient := redis.NewClient(opts)
	defer redisClient.Close()

	if err = redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalln(err)
	}

	cacheClient := cache.New(redisClient)
	rateLimiter := ratelimit.New(redisClient, 60*time.Second)

	proxyService := proxy.New(cfg.AuthServiceURL, cfg.BooksServiceURL, cfg.ServiceKey)
	handlers := handler.NewHandler(proxyService, cacheClient, time.Minute, cfg.Version)

	r := chi.NewRouter()

	r.Use(middleware.Cors)
	r.Use(middleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logging)
	r.Use(middleware.RateLimitMiddleware(&middleware.RateLimitConfig{
		Limiter:         rateLimiter,
		AnonymousLimit:  100,
		AuthorizedLimit: 1000,
		UploadLimit:     10,
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
		lgr.Info("server listening", zap.String("addr", cfg.Port))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	lgr.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		lgr.Error("graceful shutdown failed", zap.Error(err))
	} else {
		lgr.Info("graceful shutdown complete")
	}
}

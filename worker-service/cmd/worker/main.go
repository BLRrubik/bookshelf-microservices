package main

import (
	"bookshelf/worker-service/internal/config"
	"log/slog"
)

func main() {
	slog.Info("worker-service starting...")

	slog.Info("Loading configuration...")
	cfg := config.Load()

	_ = cfg
}

package config

import "os"

type Config struct {
	Port           string
	DatabaseURL    string
	AuthServiceURL string
}

func Load() *Config {
	cfg := &Config{
		Port:           os.Getenv("PORT"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AuthServiceURL: os.Getenv("AUTH_SERVICE_URL"),
	}

	normalizeConfig(cfg)

	return cfg
}

func normalizeConfig(cfg *Config) {
	if cfg.Port == "" {
		cfg.Port = ":8082"
	}

	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://postgres:postgres@localhost:5433/books?sslmode=disable"
	}

	if cfg.AuthServiceURL == "" {
		cfg.AuthServiceURL = "http://localhost:8081"
	}
}

package config

import "os"

type Config struct {
	Port                string
	DatabaseURL         string
	AuthServiceURL      string
	ServiceKey          string
	MinIOEndpoint       string // localhost:9000
	MinIOAccessKey      string
	MinIOSecretKey      string
	MinioBucket         string // bookshelf-covers
	MinIOPublicEndpoint string // http://localhost:9000
	MinIOUseSSL         bool   // false для локальной разработки
}

func Load() *Config {
	cfg := &Config{
		Port:           os.Getenv("PORT"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AuthServiceURL: os.Getenv("AUTH_SERVICE_URL"),
		ServiceKey:     os.Getenv("SERVICE_KEY"),
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

	if cfg.ServiceKey == "" {
		cfg.ServiceKey = "serviceKey"
	}

	if cfg.MinIOEndpoint == "" {
		cfg.MinIOEndpoint = "amqp://guest:guest@localhost:5672/"
	}

	if cfg.MinIOPublicEndpoint == "" {
		cfg.MinIOPublicEndpoint = "amqp://guest:guest@localhost:5672/"
	}

	if cfg.MinioBucket == "" {
		cfg.MinioBucket = "bookshelf-covers"
	}

	if cfg.MinIOAccessKey == "" {
		cfg.MinIOAccessKey = "minioadmin"
	}

	if cfg.MinIOSecretKey == "" {
		cfg.MinIOSecretKey = "minioadmin"
	}
}

package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func Load() *Config {
	cfg := &Config{
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	normalizeConfig(cfg)

	return cfg
}

func normalizeConfig(cfg *Config) {
	if cfg.Port == "" {
		cfg.Port = ":8081"
	}

	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "secret"
	}
}

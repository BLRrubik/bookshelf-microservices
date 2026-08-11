package config

import "os"

type Config struct {
	RabbitMQURL         string // URL для подключения к RabbitMQ
	DatabaseURL         string // URL для подключения к PostgreSQL (books_db)
	MinIOEndpoint       string // localhost:9000
	MinIOAccessKey      string
	MinIOSecretKey      string
	MinioBucket         string // bookshelf-covers
	MinIOPublicEndpoint string // http://localhost:9000
	MinIOUseSSL         bool   // false для локальной разработки
}

func Load() *Config {
	cfg := &Config{
		RabbitMQURL:         os.Getenv("RABBITMQ_URL"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		MinIOEndpoint:       os.Getenv("MINIO_ENDPOINT"),
		MinioBucket:         os.Getenv("MINIO_BUCKET"),
		MinIOPublicEndpoint: os.Getenv("MINIO_PUBLIC_ENDPOINT"),
		MinIOUseSSL:         os.Getenv("MINIO_USE_SSL") == "true",
		MinIOAccessKey:      os.Getenv("MINIO_ACCESS_KEY"),
		MinIOSecretKey:      os.Getenv("MINIO_SECRET_KEY"),
	}

	cfg.normalize()

	return cfg
}

func (c *Config) normalize() {
	if c.RabbitMQURL == "" {
		c.RabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	if c.MinIOEndpoint == "" {
		c.MinIOEndpoint = "amqp://guest:guest@localhost:5672/"
	}

	if c.MinIOPublicEndpoint == "" {
		c.MinIOPublicEndpoint = "amqp://guest:guest@localhost:5672/"
	}

	if c.MinioBucket == "" {
		c.MinioBucket = "bookshelf-covers"
	}

	if c.MinIOAccessKey == "" {
		c.MinIOAccessKey = "minioadmin"
	}

	if c.MinIOSecretKey == "" {
		c.MinIOSecretKey = "minioadmin"
	}

	if c.DatabaseURL == "" {
		c.DatabaseURL = "postgres://postgres:postgres@localhost:5433/books?sslmode=disable"
	}
}

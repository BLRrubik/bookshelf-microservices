package config

import "os"

type Config struct {
	RabbitMQURL string // URL для подключения к RabbitMQ
	MinIOURL    string // URL для подключения к MinIO
	MinIOAccess string // Access key для MinIO
	MinIOSecret string // Secret key для MinIO
	MinioBucket string // Имя bucket в MinIO
	DatabaseURL string // URL для подключения к PostgreSQL (books_db)
}

func Load() *Config {
	cfg := &Config{
		RabbitMQURL: os.Getenv("RABBITMQ_URL"),
		MinIOURL:    os.Getenv("MINIO_URL"),
		MinIOAccess: os.Getenv("MINIO_ACCESS"),
		MinIOSecret: os.Getenv("MINIO_SECRET"),
		MinioBucket: os.Getenv("MINIO_BUCKET"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	cfg.normalize()

	return cfg
}

func (c *Config) normalize() {
	if c.RabbitMQURL == "" {
		c.RabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	if c.MinIOURL == "" {
		c.MinIOURL = "amqp://guest:guest@localhost:5672/"
	}

	if c.MinIOAccess == "" {
		c.MinIOAccess = "read"
	}

	if c.MinIOSecret == "" {
		c.MinIOSecret = "secret"
	}

	if c.DatabaseURL == "" {
		c.DatabaseURL = os.Getenv("DATABASE_URL")
	}
}

package config

import "os"

type Config struct {
	Port                string
	RabbitMQURL         string // URL для подключения к RabbitMQ
	MinIOURL            string // URL для подключения к MinIO
	MinIOAccess         string // Access key для MinIO
	MinIOSecret         string // Secret key для MinIO
	MinioBucket         string // Имя bucket в MinIO
	MinioPublicEndpoint string
	MinioUseSSL         bool
	DatabaseURL         string // URL для подключения к PostgreSQL (books_db)
}

func Load() *Config {
	cfg := &Config{
		Port:                os.Getenv("PORT"),
		RabbitMQURL:         os.Getenv("RABBITMQ_URL"),
		MinIOURL:            os.Getenv("MINIO_URL"),
		MinioPublicEndpoint: os.Getenv("MINIO_PUBLIC_ENDPOINT"),
		MinioUseSSL:         os.Getenv("MINIO_USE_SSL") == "true",
		MinIOAccess:         os.Getenv("MINIO_ACCESS"),
		MinIOSecret:         os.Getenv("MINIO_SECRET"),
		MinioBucket:         os.Getenv("MINIO_BUCKET"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
	}

	cfg.normalize()

	return cfg
}

func (c *Config) normalize() {
	if c.Port == "" {
		c.Port = ":8083"
	}

	if c.RabbitMQURL == "" {
		c.RabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	if c.MinIOURL == "" {
		c.MinIOURL = "localhost:9000"
	}

	if c.MinioPublicEndpoint == "" {
		c.MinioPublicEndpoint = "http://localhost:9000"
	}

	if c.MinioBucket == "" {
		c.MinioBucket = "bookshelf-covers"
	}

	if c.MinIOAccess == "" {
		c.MinIOAccess = "minioadmin"
	}

	if c.MinIOSecret == "" {
		c.MinIOSecret = "minioadmin"
	}

	if c.DatabaseURL == "" {
		c.DatabaseURL = "postgres://postgres:postgres@localhost:5433/books?sslmode=disable"
	}
}

package config

import "os"

type Config struct {
	Port            string
	AuthServiceURL  string
	BooksServiceURL string
	LogFormat       string
	RedisURL        string
}

func Load() *Config {
	config := &Config{
		Port:            os.Getenv("PORT"),
		AuthServiceURL:  os.Getenv("AUTH_SERVICE_URL"),
		BooksServiceURL: os.Getenv("BOOKS_SERVICE_URL"),
		LogFormat:       os.Getenv("LOG_FORMAT"),
		RedisURL:        os.Getenv("REDIS_URL"),
	}

	config.normalize()

	return config
}

func (c *Config) normalize() {
	if c.Port == "" {
		c.Port = ":8000"
	}

	if c.AuthServiceURL == "" {
		c.AuthServiceURL = "http://localhost:8081"
	}

	if c.BooksServiceURL == "" {
		c.BooksServiceURL = "http://localhost:8082"
	}

	if c.LogFormat != "json" && c.LogFormat != "text" {
		c.LogFormat = "text"
	}

	if c.RedisURL == "" {
		c.RedisURL = "localhost:6379"
	}
}

package config

import "os"

type Config struct {
	Port            string
	AuthServiceURL  string
	BooksServiceURL string
}

func Load() *Config {
	config := &Config{
		Port:            os.Getenv("PORT"),
		AuthServiceURL:  os.Getenv("AUTH_SERVICE_URL"),
		BooksServiceURL: os.Getenv("BOOKS_SERVICE_URL"),
	}

	config.normalize()

	return config
}

func (c *Config) normalize() {
	if c.Port == "" {
		c.Port = ":8080"
	}

	if c.AuthServiceURL == "" {
		c.AuthServiceURL = "http://localhost:8081"
	}

	if c.BooksServiceURL == "" {
		c.BooksServiceURL = "http://localhost:8081"
	}
}

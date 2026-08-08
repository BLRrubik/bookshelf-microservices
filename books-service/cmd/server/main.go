package main

import (
	"bookshelf/books-service/internal/config"
	"bookshelf/books-service/internal/handler"
	"bookshelf/books-service/internal/repository"
	"bookshelf/books-service/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.Load()

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repos := repository.New(db)
	services := service.New(repos, cfg.JWTSecret)
	handlers := handler.New(services, cfg.JWTSecret)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	registerRoutes(r, handlers)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         cfg.Port,
		Handler:      r,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(r *chi.Mux, handlers *handler.Handler) {
	r.Get("/health", handlers.Health)
	r.Get("/ready", handlers.Ready)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/books", handlers.ListBooks)
		r.Get("/books/{bookId}", handlers.GetBook)
		r.Get("/books/{bookId}/reviews", handlers.ListBookReviews)

		r.Get("/reviews/{reviewId} ", handlers.GetReview)

		// Защищённые роуты — с AuthMiddleware
		r.Group(func(r chi.Router) {
			//r.Use(handlers.AuthMiddleware)

			r.Post("/books", handlers.CreateBook)
			r.Put("/books/{bookId}", handlers.UpdateBook)
			r.Delete("/books/{bookId}", handlers.DeleteBook)
			r.Post("/books/{bookId}/reviews", handlers.CreateReview)

			r.Put("/reviews/{reviewId} ", handlers.UpdateReview)
			r.Delete("/reviews/{reviewId} ", handlers.DeleteReview)
		})
	})
}

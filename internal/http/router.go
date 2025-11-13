package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/joeynolan/go-http-server/internal/http/handlers"
)

func Register(r chi.Router) {
	r.Get("/health", handlers.HealthHandler)
	r.Get("/v1/r/{code}", handlers.RedirectHandler)
	r.Post("/v1/shorten", handlers.ShortenHandler)
}

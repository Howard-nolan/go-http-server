package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/joeynolan/go-http-server/internal/http/handlers"
)

func Register(r chi.Router, h *handlers.Handler) {
	r.Get("/health", handlers.HealthHandler)
	r.Get("/readyz", h.ReadyHandler)
	r.Get("/r/{code}", h.RedirectHandler)
	r.Post("/shorten", h.ShortenHandler)
	r.Handle("/metrics", MetricsHandler())
}

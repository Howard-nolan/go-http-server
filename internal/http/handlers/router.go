package handlers

import "github.com/go-chi/chi/v5"

func Register(r chi.Router) {
	r.Get("/health", healthHandler)
	r.Get("/v1/r/{code}", redirectHandler)
	r.Post("/v1/shorten", shortenHandler)
}

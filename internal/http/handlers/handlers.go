package handlers

import (
	"database/sql"

	lru "github.com/hashicorp/golang-lru/v2"
	ilog "github.com/joeynolan/go-http-server/internal/platform/log"
)

type Handler struct {
	DB    *sql.DB
	Log   *ilog.Logger
	cache *lru.Cache[string, string]
}

func NewHandler(db *sql.DB, logger *ilog.Logger) *Handler {
	cache, _ := lru.New[string, string](32)
	return &Handler{
		DB:    db,
		Log:   logger,
		cache: cache,
	}
}

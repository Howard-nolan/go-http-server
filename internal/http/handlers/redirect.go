package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	ilog "github.com/joeynolan/go-http-server/internal/platform/log"
)

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	logger := ilog.New()

	code := chi.URLParam(r, "code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "missing code")
		return
	}

	var url string
	logger.Infof("Code: %s", code)
	err := h.DB.QueryRowContext(r.Context(), "SELECT url FROM links WHERE code = ?", code).Scan(&url)
	logger.Infof("URL: %s", url)
	if err == sql.ErrNoRows {
		WriteError(w, http.StatusNotFound, "code not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "lookup failed")
		return
	}

	target := url
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "https://" + target
	}
	http.Redirect(w, r, target, http.StatusFound)
}

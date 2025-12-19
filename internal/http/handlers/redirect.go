package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "missing code")
		return
	}

	var url string
	h.Log.Infof("Code: %s", code)
	if val, ok := h.cache.Get(code); ok {
		h.Log.Infof("Cache hit for code: %s", code)
		url = val
	} else {
		err := h.DB.QueryRowContext(r.Context(), "SELECT url FROM links WHERE code = ?", code).Scan(&url)
		if err == sql.ErrNoRows {
			WriteError(w, http.StatusNotFound, "code not found")
			return
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			WriteError(w, http.StatusRequestTimeout, "timeout")
			return
		}
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "lookup failed")
			return
		}
		h.cache.Add(code, url)
	}
	h.Log.Infof("URL: %s", url)

	target := url
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "https://" + target
	}
	http.Redirect(w, r, target, http.StatusFound)
}

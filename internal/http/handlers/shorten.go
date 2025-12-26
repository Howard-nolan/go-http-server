package handlers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {

	var req struct {
		URL    string `json:"url"`
		Expiry int    `json:"expiry,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.URL == "" {
		WriteError(w, http.StatusBadRequest, "url is required")
		return
	}

	var requestID string = middleware.GetReqID(r.Context())
	// Generate a random code; insert with the URL; retry on collision
	var code string
	for attempts := 0; attempts < 3; attempts++ {
		var err error
		code, err = randomCode(6)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to generate code")
			return
		}

		res, err := h.DB.ExecContext(r.Context(), "INSERT INTO links (code, url, idempotency_key) VALUES (?, ?, ?) ON CONFLICT (idempotency_key) DO NOTHING", code, req.URL, requestID)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			WriteError(w, http.StatusRequestTimeout, "timeout")
			return
		}

		if err != nil {
			if attempts == 2 {
				WriteError(w, http.StatusInternalServerError, "failed to generate unique code")
				return
			}
			continue
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			if err := h.DB.QueryRowContext(r.Context(), "SELECT code FROM links WHERE idempotency_key = ?", requestID).Scan(&code); err != nil {
				WriteError(w, http.StatusInternalServerError, "DB Access Error for Duplicate Request")
				return
			}
		}
		break
	}

	resp := map[string]string{
		"short": fmt.Sprintf("%s/%s", h.BaseURL, code), // generated value
	}
	WriteJSON(w, http.StatusCreated, resp)
}

func randomCode(length int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		b[i] = letters[int(b[i])%len(letters)]
	}

	return string(b), nil
}

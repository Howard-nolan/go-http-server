package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
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

	// Generate a random code; insert with the URL; retry on collision
	var code string
	for attempts := 0; attempts < 3; attempts++ {
		var err error
		code, err = randomCode(6)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to generate code")
			return
		}

		_, err = h.DB.ExecContext(r.Context(), "INSERT INTO links (code, url) VALUES (?, ?)", code, req.URL)
		if err == nil {
			break
		}

		// throw an error if max attempts reached
		if attempts == 2 {
			WriteError(w, http.StatusInternalServerError, "failed to generate unique code")
			return
		}
	}

	resp := map[string]string{
		"short": fmt.Sprintf("https://short.example/%s", code), // generated value
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

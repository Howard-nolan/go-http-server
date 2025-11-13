package handlers

import (
	"encoding/json"
	"net/http"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {

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

	resp := map[string]string{
		"short": "https://short.example/abc123", // generated value
	}
	WriteJSON(w, http.StatusCreated, resp)
}

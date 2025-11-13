package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is a simple JSON error body with a human-friendly message.
type ErrorResponse struct {
	Message string `json:"message"`
}

// WriteJSON writes v as a JSON response with the provided status code and
// sets the Content-Type to application/json; charset=utf-8.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError writes a standardized JSON error response containing only a
// status code and a message. Example body: {"message":"not found"}
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Message: message})
}

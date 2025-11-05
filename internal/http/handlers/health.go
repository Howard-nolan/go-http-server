package handlers

import (
	"encoding/json"
	"net/http"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", healthHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

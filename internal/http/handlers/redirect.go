package handlers

import (
	"encoding/json"
	"net/http"
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"status": "ok", "message": "redirect link"}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

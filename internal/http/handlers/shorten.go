package handlers

import (
	"encoding/json"
	"net/http"
)

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"status": "ok", "message": "shorten endpoint"}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

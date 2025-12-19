package handlers

import (
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"status": "ok"}
	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.PingContext(r.Context()); err != nil {
		WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	resp := map[string]string{"status": "ok"}
	WriteJSON(w, http.StatusOK, resp)
}

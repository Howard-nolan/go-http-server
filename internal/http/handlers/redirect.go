package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {

	code := chi.URLParam(r, "code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "missing code")
		return
	}

	resp := map[string]string{"message": "redirect link"}
	WriteJSON(w, http.StatusFound, resp)
}

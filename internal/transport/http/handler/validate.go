package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (hd *handlerStr) validate(rt *chi.Mux) {
	header := r.Header.Get("Authorization")

	if header == "" {
		http.Error(w, "Заголовок авторизации отсутствует", http.StatusUnauthorized)
		return
	}
}

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golovanevvs/gophermart/internal/model"
)

func (hd *HandlerStr) Register(w http.ResponseWriter, r *http.Request) {
	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// десериализация JSON
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Установка Content-Type
	w.Header().Set("Content-Type", "text/plain")

	// установка заголовков
	w.WriteHeader(http.StatusOK)

	// запись ответа
	w.Write([]byte(user.Login))
	w.Write([]byte(user.Password))
}

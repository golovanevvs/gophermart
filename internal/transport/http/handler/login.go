package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golovanevvs/gophermart/internal/model"
)

func (hd *handlerStr) login(w http.ResponseWriter, r *http.Request) {
	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат запроса"))
		return
	}

	// десериализация JSON
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte("Внутренняя ошибка сервера"))
		return
	}
}

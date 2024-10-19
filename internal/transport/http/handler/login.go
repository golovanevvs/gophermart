package handler

import (
	"encoding/json"
	"fmt"
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

	// получение строки токена
	tokenString, customErr := hd.sv.BuildJWTString(r.Context(), user.Login, user.Password)
	if customErr.IsError {
		http.Error(w, customErr.CustomErr.Error(), http.StatusInternalServerError)
		return
	}

	// отправка заголовков
	w.Header().Set("Authorization", fmt.Sprint("Bearer", tokenString))
	w.WriteHeader(http.StatusOK)
}

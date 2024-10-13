package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golovanevvs/gophermart/internal/model"
)

func (hd *handlerStr) register(w http.ResponseWriter, r *http.Request) {
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

	err := hd.sv.AuthInt.CreateUser(r.Context(), user)
	if err != nil {
		// если логин уже существует в БД
		if strings.Contains(err.Error(), "Unique") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Логин уже занят"))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte("Внутренняя ошибка сервера"))
		return
	}

	// запись заголовков и ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно зарегистрирован"))
}

package handler

import (
	"encoding/json"
	"fmt"
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

	// запуск сервиса CreateUser, проверка ошибок
	userID, err := hd.sv.AuthInt.CreateUser(r.Context(), user)
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
	user.UserID = userID

	// запись заголовков и ответа
	w.WriteHeader(http.StatusOK)
	res := fmt.Sprintf("Пользователь %v успешно зарегистрирован под номером %v", user.Login, user.UserID)
	w.Write([]byte(res))
}

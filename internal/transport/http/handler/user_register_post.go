package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golovanevvs/gophermart/internal/model"
)

func (hd *handlerStr) userRegister(w http.ResponseWriter, r *http.Request) {
	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат запроса"))
		return
	}

	// десериализация JSON (user.Login, user.Password)
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("Ошибка в NewDecoder")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte("Внутренняя ошибка сервера"))
		return
	}

	// запуск сервиса CreateUser, проверка ошибок
	userID, err := hd.sv.AuthServiceInt.CreateUser(r.Context(), user)
	if err != nil {
		// если логин уже существует в БД
		if strings.Contains(err.Error(), "Unique") {
			fmt.Println("Ошибка в CreateUser1")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Логин уже занят"))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Ошибка в CreateUser2")
		w.Write([]byte("Внутренняя ошибка сервера"))
		return
	}
	user.UserID = userID

	// получение строки токена
	tokenString, customErr := hd.sv.BuildJWTString(r.Context(), user.Login, user.Password)
	if customErr.IsError {
		fmt.Println("Ошибка в BuildJWTString")
		fmt.Println(customErr.CustomErr.Error())
		fmt.Println(customErr.Err)
		http.Error(w, customErr.CustomErr.Error(), http.StatusInternalServerError)
		return
	}

	// запись заголовков и ответа
	w.Header().Set("Authorization", fmt.Sprint("Bearer", tokenString))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	res := fmt.Sprintf("Пользователь %v успешно зарегистрирован под номером %v", user.Login, user.UserID)
	w.Write([]byte(res))
	w.Write([]byte(tokenString))
}

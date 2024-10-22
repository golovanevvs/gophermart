package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (hd *handlerStr) userRegister(w http.ResponseWriter, r *http.Request) {
	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Ошибка %v: %v. Требуется: application/json. Получено: %v", http.StatusBadRequest, string(customerrors.InvalidContentType400), contentType)))
		return
	}

	// десериализация JSON (user.Login, user.Password)
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Ошибка %v: %v, %v", http.StatusInternalServerError, string(customerrors.ErrorDecodeJSON500), err.Error())))
		return
	}

	// запуск сервиса CreateUser, проверка ошибок
	// TODO: Обработать кастомные ошибки
	userID, err := hd.sv.AuthServiceInt.CreateUser(r.Context(), user)
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

	// получение строки токена
	tokenString, customErr := hd.sv.BuildJWTString(r.Context(), user.Login, user.Password)
	if customErr.IsError {
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

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (hd *handlerStr) userRegister(w http.ResponseWriter, r *http.Request) {
	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
	default:
		http.Error(w, string(customerrors.InvalidContentType400), http.StatusBadRequest)
		return
	}

	// десериализация JSON в user (user.Login, user.Password)
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, string(customerrors.DecodeJSONError500), http.StatusInternalServerError)
		return
	}

	// запуск сервиса CreateUser, получение userID нового пользователя для последующей авторизации, проверка ошибок
	userID, customErr := hd.sv.AuthServiceInt.CreateUser(r.Context(), user)
	if customErr.IsError {
		switch customErr.CustomErr {
		case customerrors.DBBusyLogin409:
			// если логин уже существует в БД
			http.Error(w, customErr.AllErr.Error(), http.StatusConflict)
			return
			// прочие ошибки
		case customerrors.DBError500:
			http.Error(w, customErr.AllErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	// сохранение userID в user, если не было ошибок (user.UserID)
	user.UserID = userID

	// авторизация
	// получение строки токена
	tokenString, customErr := hd.sv.BuildJWTString(r.Context(), user.Login, user.Password)
	if customErr.IsError {
		http.Error(w, customErr.AllErr.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]interface{})
	resMap["Login"] = user.Login
	resMap["userID"] = user.UserID
	resMap["token"] = tokenString

	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// запись заголовков и ответа
	w.Header().Add("Authorization", fmt.Sprint("Bearer", tokenString))
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	//res := fmt.Sprintf("Пользователь %v успешно зарегистрирован под номером %v\n", user.Login, user.UserID)
	w.Write([]byte(res))
	//w.Write([]byte(tokenString))
}

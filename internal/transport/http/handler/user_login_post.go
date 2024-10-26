package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (hd *handlerStr) userLogin(w http.ResponseWriter, r *http.Request) {
	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
	default:
		http.Error(w, string(customerrors.InvalidContentType400), http.StatusBadRequest)
		return
	}

	// десериализация JSON
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, string(customerrors.DecodeJSONError500), http.StatusInternalServerError)
		return
	}

	// получение строки токена
	tokenString, customErr := hd.sv.BuildJWTString(r.Context(), user.Login, user.Password)
	if customErr.IsError {
		switch customErr.CustomErr {
		case customerrors.DBInvalidLoginPassword401:
			http.Error(w, customErr.AllErr.Error(), http.StatusUnauthorized)
			return
		case customerrors.InternalServerError500:
			http.Error(w, customErr.AllErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	// формирование ответа
	resMap := make(map[string]interface{})
	resMap["Login"] = user.Login
	resMap["userID"] = user.UserID
	resMap["token"] = tokenString
	res, err := json.MarshalIndent(resMap, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// отправка заголовков и ответа
	w.Header().Set("Authorization", fmt.Sprint("Bearer", tokenString))
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

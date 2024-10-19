package handler

import (
	"fmt"
	"net/http"
)

func (hd *handlerStr) orders(w http.ResponseWriter, r *http.Request) {
	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат запроса"))
		return
	}

	// проверка, что пользователь авторизован
	userID := r.Context().Value("userIDContextKey")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte(fmt.Sprintf("Orders: userID = %v", userID)))
}

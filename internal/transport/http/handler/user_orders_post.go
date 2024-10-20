package handler

import (
	"fmt"
	"net/http"
)

func (hd *handlerStr) userUploadOrder(w http.ResponseWriter, r *http.Request) {
	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "text/plain":
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат запроса"))
		return
	}

	// получение userID
	userID := r.Context().Value(UserIDContextKey)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte(fmt.Sprintf("Orders: userID = %v", userID)))
}

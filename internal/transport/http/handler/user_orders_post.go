package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	userID := r.Context().Value(UserIDContextKey).(int)

	// чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат запроса"))
		return
	}

	// получение номера заказа
	orderNumber, err := strconv.Atoi(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат запроса"))
		return
	}

	// запуск сервиса
	orderID, err := hd.sv.UploadOrder(r.Context(), userID, orderNumber)
	// TODO: добавить обработку кастомных ошибок
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат запроса"))
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprintf("Заказ от пользователя с userID %v принят под номером %v", userID, orderID)))
}

package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/golovanevvs/gophermart/internal/customerrors"
)

func (hd *handlerStr) userUploadOrder(w http.ResponseWriter, r *http.Request) {
	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "text/plain":
	default:
		http.Error(w, string(customerrors.InvalidContentType400), http.StatusBadRequest)
		return
	}

	// получение userID
	userID := r.Context().Value(UserIDContextKey).(int)

	// чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, string(customerrors.InvalidRequest400), http.StatusBadRequest)
		return
	}

	// получение номера заказа
	orderNumber, err := strconv.Atoi(string(body))
	if err != nil {
		http.Error(w, string(customerrors.InvalidRequest400), http.StatusBadRequest)
		return
	}

	// запуск сервиса и обработка ошибок
	orderID, customErr := hd.sv.UploadOrder(r.Context(), userID, orderNumber)
	if customErr.IsError {
		switch customErr.CustomErr {
		case customerrors.OrderAlredyUploadedThisUser200:
			http.Error(w, customErr.AllErr.Error(), http.StatusOK)
			return
		case customerrors.OrderAlredyUploadedOtherUser409:
			http.Error(w, customErr.AllErr.Error(), http.StatusConflict)
			return
		case customerrors.DBError500:
			http.Error(w, customErr.AllErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	//запись заголовков и ответа
	w.WriteHeader(http.StatusAccepted)

	w.Write([]byte(fmt.Sprintf("Заказ от пользователя с userID %v принят под номером %v", userID, orderID)))
}

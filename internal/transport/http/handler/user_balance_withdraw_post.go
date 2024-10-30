package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (hd *handlerStr) withDraw(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать withDraw

	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
	default:
		http.Error(w, string(customerrors.InvalidContentType400), http.StatusBadRequest)
		return
	}

	// получение userID
	userID := r.Context().Value(UserIDContextKey).(int)

	// десериализация JSON в model.WithdrawOrder
	var withdrawOrder model.WithdrawOrder
	if err := json.NewDecoder(r.Body).Decode(&withdrawOrder); err != nil {
		http.Error(w, string(customerrors.DecodeJSONError500), http.StatusInternalServerError)
		return
	}

	// запуск сервиса и обработка ошибок
	customError := hd.sv.Withdraw(r.Context(), userID, withdrawOrder.OrderNumber, withdrawOrder.Sum)
	if customError.IsError {
		switch customError.CustomErr {
		case customerrors.InvalidOrderNumberNotInt422:
			http.Error(w, customError.AllErr.Error(), http.StatusUnprocessableEntity)
			return
		case customerrors.InvalidOrderNumber422:
			http.Error(w, customError.AllErr.Error(), http.StatusUnprocessableEntity)
			return
		case customerrors.DBError500:
			http.Error(w, customError.AllErr.Error(), http.StatusInternalServerError)
			return
		case customerrors.NotEnoughPoints402:
			http.Error(w, customError.AllErr.Error(), http.StatusPaymentRequired)
			return
		}
	}

	// запись заголовков и ответа
	w.WriteHeader(http.StatusOK)

	// TODO: сделать нормальный ответ
	w.Write([]byte("withDraw"))
}

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (hd *handlerStr) withDraw(w http.ResponseWriter, r *http.Request) {
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
	var withdrawOrder model.Withdrawals
	if err := json.NewDecoder(r.Body).Decode(&withdrawOrder); err != nil {
		http.Error(w, fmt.Errorf("%v: %v", customerrors.DecodeJSONError500, err).Error(), http.StatusInternalServerError)
		return
	}

	// запуск сервиса и обработка ошибок
	err := hd.sv.Withdraw(r.Context(), userID, withdrawOrder.NewOrderNumber, withdrawOrder.Sum)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), customerrors.InvalidOrderNumberNotInt422):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		case strings.Contains(err.Error(), customerrors.InvalidOrderNumber422):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		case strings.Contains(err.Error(), customerrors.DBError500):
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		case strings.Contains(err.Error(), customerrors.NotEnoughPoints402):
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}
	}

	// запись заголовков и ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("успешная обработка запроса"))
}

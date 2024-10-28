package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golovanevvs/gophermart/internal/customerrors"
)

func (hd *handlerStr) getOrders(w http.ResponseWriter, r *http.Request) {
	// TODO: добавить accrual
	// TODO: добавить обновление статуса

	// получение userID
	userID := r.Context().Value(UserIDContextKey).(int)

	// запуск сервиса и обработка ошибок
	orders, customErr := hd.sv.GetOrders(r.Context(), userID)
	if customErr.IsError {
		switch customErr.CustomErr {
		// если нет загруженных заказов
		case customerrors.EmptyOrder204:
			http.Error(w, customErr.AllErr.Error(), http.StatusNoContent)
			return
		// прочие ошибки
		case customerrors.DBError500:
			http.Error(w, customErr.AllErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	// формирование тела ответа
	res, err := json.MarshalIndent(orders, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// запись заголовков и ответа
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golovanevvs/gophermart/internal/customerrors"
)

func (hd *handlerStr) getOrders(w http.ResponseWriter, r *http.Request) {
	// TODO: добавить accrual
	// TODO: добавить обновление статуса

	// получение userID
	userID := r.Context().Value(UserIDContextKey).(int)

	// запуск сервиса и обработка ошибок
	orders, err := hd.sv.GetOrders(r.Context(), userID)
	if err != nil {
		switch {
		// если нет загруженных заказов
		case strings.Contains(err.Error(), customerrors.EmptyOrder204):
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		// прочие ошибки
		case strings.Contains(err.Error(), customerrors.DBError500):
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// формирование тела ответа
	res, err := json.MarshalIndent(orders, "", " ")
	if err != nil {
		http.Error(w, fmt.Errorf("%v: %v", customerrors.EncodeJSONError500, err.Error()).Error(), http.StatusInternalServerError)
		return
	}

	// запись заголовков и ответа
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

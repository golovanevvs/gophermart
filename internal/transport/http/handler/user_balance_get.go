package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/gophermart/internal/customerrors"
)

func (hd *handlerStr) getBalance(w http.ResponseWriter, r *http.Request) {
	// получение userID
	userID := r.Context().Value(UserIDContextKey).(int)

	// запуск сервиса и обработка ошибок
	balance, err := hd.sv.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// формирование тела ответа
	res, err := json.MarshalIndent(balance, "", " ")
	if err != nil {
		http.Error(w, fmt.Errorf("%v: %v", customerrors.EncodeJSONError500, err.Error()).Error(), http.StatusInternalServerError)
		return
	}

	// запись заголовков и ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

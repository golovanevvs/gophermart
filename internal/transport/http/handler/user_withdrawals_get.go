package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golovanevvs/gophermart/internal/customerrors"
)

func (hd *handlerStr) withDrawals(w http.ResponseWriter, r *http.Request) {
	// получение userID
	userID := r.Context().Value(UserIDContextKey).(int)

	// запуск сервиса и обработка ошибок
	withdrawals, err := hd.sv.Withdrawals(r.Context(), userID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), customerrors.EmptyWithdrawals204):
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		case strings.Contains(err.Error(), customerrors.DBError500):
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// формирование тела ответа
	res, err := json.MarshalIndent(withdrawals, "", " ")
	if err != nil {
		http.Error(w, fmt.Errorf("%v: %v", customerrors.EncodeJSONError500, err.Error()).Error(), http.StatusInternalServerError)
		return
	}

	// запись заголовков и ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

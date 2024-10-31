package handler

import (
	"encoding/json"
	"net/http"
)

func (hd *handlerStr) withDrawals(w http.ResponseWriter, r *http.Request) {
	// получение userID
	userID := r.Context().Value(UserIDContextKey).(int)

	// запуск сервиса и обработка ошибок
	withdrawals, err := hd.sv.Withdrawals(r.Context(), userID)
	if err != nil {
		// TODO: добавить обработку ошибок
	}

	// формирование тела ответа
	res, err := json.MarshalIndent(withdrawals, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// запись заголовков и ответа
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

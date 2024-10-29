package handler

import (
	"encoding/json"
	"net/http"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// запись заголовков и ответа
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

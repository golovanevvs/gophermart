package handler

import (
	"encoding/json"
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
		//
		http.Error(w, string(customerrors.InvalidOrderNumber422), http.StatusUnprocessableEntity)
		return
	}

	// запуск сервиса и обработка ошибок
	orderID, customErr := hd.sv.UploadOrder(r.Context(), userID, orderNumber)
	if customErr.IsError {
		switch customErr.CustomErr {
		// номер заказа не соответствует алгоритму Луна
		case customerrors.InvalidOrderNumber422:
			http.Error(w, customErr.AllErr.Error(), http.StatusUnprocessableEntity)
		// номер заказа уже был загружен этим пользователем
		case customerrors.OrderAlredyUploadedThisUser200:
			http.Error(w, customErr.AllErr.Error(), http.StatusOK)
			return
		// номер заказа уже был загружен другим пользователем
		case customerrors.OrderAlredyUploadedOtherUser409:
			http.Error(w, customErr.AllErr.Error(), http.StatusConflict)
			return
		// внутренняя ошибка сервера
		case customerrors.DBError500:
			http.Error(w, customErr.AllErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	// формирование ответа
	resMap := make(map[string]interface{})
	resMap["userID"] = userID
	resMap["orderID"] = orderID
	resMap["orderNumber"] = orderNumber

	res, err := json.MarshalIndent(resMap, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//запись заголовков и ответа
	w.WriteHeader(http.StatusAccepted)
	w.Write(res)
}

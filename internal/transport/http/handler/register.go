package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golovanevvs/gophermart/internal/model"
)

func (hd *handlerStr) register(w http.ResponseWriter, r *http.Request) {
	// проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// десериализация JSON
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := hd.sv.AuthInt.CreateUser(r.Context(), user)
	if err != nil {
		// если логин уже существует в БД
		if strings.Contains(err.Error(), "Unique") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Такой login уже существует"))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// установка заголовков
	w.WriteHeader(http.StatusOK)
}

package handler

import "net/http"

func (hd *handlerStr) getBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("getBalance"))
}
package handler

import "net/http"

func (hd *handlerStr) getBalance(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать getBalance
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("getBalance"))
}

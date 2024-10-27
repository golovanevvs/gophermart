package handler

import "net/http"

func (hd *handlerStr) getOrderNumbers(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать getOrderNumber
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("getOrderNumber"))
}

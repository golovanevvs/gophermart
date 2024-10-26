package handler

import "net/http"

func (hd *handlerStr) withDraw(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать withDraw
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("withDraw"))
}

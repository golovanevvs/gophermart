package handler

import "net/http"

func (hd *handlerStr) withDrawals(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать withDrawals
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("withDrawals"))
}

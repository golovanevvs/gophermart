package handler

import (
	"net/http"
)

func (hd *HandlerStr) Register(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(hd.sv.CreateUser()))
	w.WriteHeader(http.StatusOK)
}

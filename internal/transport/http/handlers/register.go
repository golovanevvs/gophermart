package handlers

import "net/http"

func (hd HandlersStr) Register(w http.ResponseWriter, r *http.Request) {
	w1 := hd.service.RegisterUser()
	w.Write([]byte(w1))
}

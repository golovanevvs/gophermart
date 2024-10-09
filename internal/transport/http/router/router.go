package router

import (
	"net/http"

	"github.com/go-chi/chi"
)

type HandlersInt interface {
	Register(w http.ResponseWriter, r *http.Request)
}

func New(hd HandlersInt) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", hd.Register)
	})

	return r
}

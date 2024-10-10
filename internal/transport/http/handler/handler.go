package handler

import (
	"github.com/go-chi/chi"
	"github.com/golovanevvs/gophermart/internal/service"
)

type HandlerStr struct {
	sv *service.ServiceStr
}

func NewHandler(sv *service.ServiceStr) *HandlerStr {
	return &HandlerStr{
		sv: sv,
	}
}

func (hd *HandlerStr) InitRoutes() *chi.Mux {
	rt := chi.NewRouter()
	rt.Route("/api/user", func(r chi.Router) {
		r.Post("/register", hd.Register)
	})
	return rt
}

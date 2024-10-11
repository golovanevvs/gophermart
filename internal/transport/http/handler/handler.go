package handler

import (
	"github.com/go-chi/chi"
	"github.com/golovanevvs/gophermart/internal/middleware/logger"
	"github.com/golovanevvs/gophermart/internal/service"
	"github.com/sirupsen/logrus"
)

type HandlerStr struct {
	sv *service.ServiceStr
}

func NewHandler(sv *service.ServiceStr) *HandlerStr {
	return &HandlerStr{
		sv: sv,
	}
}

func (hd *HandlerStr) InitRoutes(lg *logrus.Logger) *chi.Mux {
	rt := chi.NewRouter()
	rt.Use(logger.WithLogging(lg))
	rt.Route("/api/user", func(r chi.Router) {
		r.Post("/register", hd.Register)
	})
	return rt
}

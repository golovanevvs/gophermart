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

// NewHandler - конструктор *HandlerStr
func NewHandler(sv *service.ServiceStr) *HandlerStr {
	return &HandlerStr{
		sv: sv,
	}
}

// InitRoutes - маршрутизация запросов
func (hd *HandlerStr) InitRoutes(lg *logrus.Logger) *chi.Mux {
	// создание экземпляра роутера
	rt := chi.NewRouter()

	// использование middleware
	// логгирование
	rt.Use(logger.WithLogging(lg))

	// маршруты
	rt.Route("/api/user", func(r chi.Router) {
		r.Post("/register", hd.Register)
	})

	return rt
}

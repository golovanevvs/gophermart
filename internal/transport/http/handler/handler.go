package handler

import (
	"github.com/go-chi/chi"
	"github.com/golovanevvs/gophermart/internal/middleware/logger"
	"github.com/golovanevvs/gophermart/internal/service"
	"github.com/sirupsen/logrus"
)

type handlerStr struct {
	sv *service.ServiceStrInt
}

// NewHandler - конструктор *HandlerStr
func NewHandler(sv *service.ServiceStrInt) *handlerStr {
	return &handlerStr{
		sv: sv,
	}
}

// InitRoutes - маршрутизация запросов, используется в качестве http.Handler при запуске сервера
func (hd *handlerStr) InitRoutes(lg *logrus.Logger) *chi.Mux {
	// создание экземпляра роутера
	rt := chi.NewRouter()

	// использование middleware
	// логгирование
	rt.Use(logger.WithLogging(lg))

	// маршруты
	rt.Route("/api/user", func(r chi.Router) {
		r.Post("/register", hd.userRegister)
		r.Post("/login", hd.userLogin)
		r.With(hd.authByJWT).Post("/orders", hd.userUploadOrder)
		r.With(hd.authByJWT).Get("/orders", hd.getOrders)
		r.With(hd.authByJWT).Get("/withdrawals", hd.withDrawals)
		r.With(hd.authByJWT).Route("/balance", func(r chi.Router) {
			r.Get("/", hd.getBalance)
			r.Post("/withdraw", hd.withDraw)
		})
	})

	return rt
}

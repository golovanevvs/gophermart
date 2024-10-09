package app

import (
	"net/http"

	"github.com/golovanevvs/gophermart/internal/database/postgres"
	"github.com/golovanevvs/gophermart/internal/service"
	"github.com/golovanevvs/gophermart/internal/transport/http/handlers"
	"github.com/golovanevvs/gophermart/internal/transport/http/router"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	lg := logrus.New()
	lg.SetLevel(logrus.DebugLevel)

	db := postgres.New()
	sv := service.New(db)
	hd := handlers.New(sv)
	rt := router.New(hd)

	lg.Infof("Сервер накопительной системы лояльности Гофермарт запущен")
	if err := http.ListenAndServe(":8080", rt); err != nil {
		lg.Errorf("Ошибка сервера: %v", err)
	}
}

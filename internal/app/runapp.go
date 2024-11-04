package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golovanevvs/gophermart/internal/config"
	"github.com/golovanevvs/gophermart/internal/service"
	"github.com/golovanevvs/gophermart/internal/storage"
	"github.com/golovanevvs/gophermart/internal/storage/postgres"

	"github.com/golovanevvs/gophermart/internal/transport/http/handler"

	"github.com/sirupsen/logrus"
)

func RunApp() {
	// инициализация логгера
	lg := logrus.New()
	lg.SetLevel(logrus.DebugLevel)

	// инициализация конфигурации
	cfg := config.NewConfig()

	// инициализация БД Postgres
	db, err := postgres.NewPostgres(cfg.DatabaseURI)
	if err != nil {
		lg.Fatalf("Ошибка инициализации базы данных: %v", err.Error())
	}

	// инициализация хранилища
	st := storage.NewStorage(db)
	// инициализация сервиса
	sv := service.NewService(st, cfg.AccrualSystemAddress)
	// инициализация хендлера
	hd := handler.NewHandler(sv)
	// инициализация сервера
	srv := NewServer()

	// запуск сервера
	go func() {
		lg.Infof("Сервер накопительной системы лояльности Гофермарт запущен")
		if err := srv.RunServer(cfg.RunAddress, hd.InitRoutes(lg)); err != nil {
			lg.Fatalf("Ошибка запуска сервера: %v", err.Error())
		}
	}()

	// завершение работы сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	lg.Infof("Получен сигнал о завершении работы сервера")

	if err := srv.ShutdownServer(context.Background()); err != nil {
		lg.Errorf("Ошибка при завершении работы сервера: %v", err.Error())
	}

	if err := db.Close(); err != nil {
		lg.Errorf("Ошибка при завершении работы с БД: %v", err.Error())
	}

	lg.Infof("Работа сервера завершена")
	time.Sleep(time.Second * 2)
}

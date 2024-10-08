package app

import (
	"github.com/go-chi/chi"
	"github.com/golovanevvs/gophermart/internal/config"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *chi.Mux
	logger *logrus.Logger
	config *config.Config
}

func StartServer() {

}

func New() {
	logger := logrus.New()
}

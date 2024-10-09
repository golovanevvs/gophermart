package app

import (
	"fmt"
	"net/http"

	"github.com/golovanevvs/gophermart/internal/database/postgres"
	"github.com/golovanevvs/gophermart/internal/services"
	"github.com/golovanevvs/gophermart/internal/transport/http/handlers"
	"github.com/golovanevvs/gophermart/internal/transport/http/router"
)

func StartServer() {
	db := postgres.New()
	sv := services.New(db)
	hd := handlers.New(sv)
	rt := router.New(hd)

	if err := http.ListenAndServe(":8080", rt); err != nil {
		fmt.Printf("Ошибка сервера: %v\n", err)
	}
}

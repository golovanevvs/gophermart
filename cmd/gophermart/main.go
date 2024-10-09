package main

import (
	"fmt"
	"time"

	"github.com/golovanevvs/gophermart/internal/app"
)

func main() {
	app.StartServer()
	fmt.Println("Завершение работы")
	time.Sleep(time.Second * 3)
}

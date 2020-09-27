package main

import (
	"avito/internal/pkg/config"
	"avito/internal/pkg/middleware"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"

	"avito/internal/pkg/balance"
	"avito/internal/pkg/database"
	"avito/internal/pkg/logger"

	_ "github.com/lib/pq"
)

func main() {

	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.NewPool(&c)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := database.NewRepo(pool)
	l := logger.NewLogger("INFO")
	bal := balance.NewHandler(repo)
	m := middleware.NewMiddleware(l)

	router := routing.New()
	api := router.Group("/api")
	api.Use(m.PanicMiddleware, m.LogMiddleware)
	api.Post("/transfer", bal.Transfer)
	api.Get("/balance/<id>", bal.GetBalance)
	api.Get("/list/<id>", bal.GetTransferList)

	panic(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}
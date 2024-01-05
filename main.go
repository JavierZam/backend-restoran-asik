package main

import (
	"log"
	"net/http"
	"restoran-asik/handler"
	"restoran-asik/repository"
	"restoran-asik/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := repository.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:5500"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	jwtSecret := "inisecretku"

	transactionRepo := repository.NewTransactionRepository(db)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	reportRepo := repository.NewReportRepository(db)

	uh := handler.NewUserHandler(*userRepo, jwtSecret)
	ph := handler.NewProductHandler(*productRepo)
	th := handler.NewTransactionHandler(*transactionRepo, *productRepo, *userRepo)
	rh := handler.NewReportHandler(*reportRepo, *productRepo, *transactionRepo)

	router.InitUserRoutes(e, uh)
	router.InitProductRoutes(e, ph)
	router.InitTransactionRoutes(e, th)
	router.InitReportRoutes(e, rh)

	log.Fatal(e.Start(":8080"))
}

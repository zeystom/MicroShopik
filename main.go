package main

import (
	"MicroShopik/internal/controllers"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	e := echo.New()

	helloCtrl := controllers.NewHelloController()

	e.GET("/", helloCtrl.SayHello)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

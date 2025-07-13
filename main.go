package main

import (
	_ "MicroShopik/docs" // This is the generated swagger docs
	"MicroShopik/internal/controllers"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title MicroShopik API
// @version 1.0
// @description This is a sample server for MicroShopik.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	e := echo.New()

	helloCtrl := controllers.NewHelloController()

	e.GET("/aboba", helloCtrl.SayHello)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

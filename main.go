package main

import (
	"MicroShopik/configs"
	_ "MicroShopik/docs"
	"MicroShopik/internal/controllers"
	"database/sql"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
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
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("error db connect: ", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("error pinging db: ", err)
	}

	e := echo.New()

	helloCtrl := controllers.NewHelloController()

	e.GET("/aboba", helloCtrl.SayHello)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

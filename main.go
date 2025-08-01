package main

import (
	"MicroShopik/configs"
	_ "MicroShopik/docs"
	"MicroShopik/internal/controllers"
	"MicroShopik/internal/database"
	"MicroShopik/internal/middleware"
	"MicroShopik/internal/repositories"
	"MicroShopik/internal/services"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize GORM database
	err = database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	e := echo.New()

	helloCtrl := controllers.NewHelloController()
	userRepo := repositories.NewUserRepository(database.GetDB())
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authCtrl := controllers.NewAuthController(authService)

	authGroup := e.Group("")
	authGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	e.GET("/aboba", helloCtrl.SayHello)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/register", authCtrl.Register)
	e.POST("/login", authCtrl.Login)

	// GetMe godoc
	// @Summary Get current user info
	// @Description Get the authenticated user's ID
	// @Tags auth
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]string
	// @Security ApiKeyAuth
	// @Router /me [get]

	e.GET("/me", func(c echo.Context) error {
		userID := c.Get("user_id").(int)
		return c.JSON(200, map[string]interface{}{"user_id": userID})
	})

	go func() {
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

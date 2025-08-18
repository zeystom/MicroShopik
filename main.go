package main

import (
	"MicroShopik/configs"
	_ "MicroShopik/docs"
	"MicroShopik/internal/container"
	"MicroShopik/internal/database"
	"MicroShopik/internal/middleware"
	"MicroShopik/scripts"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	err = database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	log.Println("Checking if seed data is needed...")
	if err := scripts.RunInitialSeed(cfg); err != nil {
		log.Printf("Warning: Failed to run initial seed: %v", err)
	}

	e := echo.New()

	setupCORS(e)

	newContainer := container.NewContainer()

	setupRoutes(e, newContainer, cfg.JWTSecret)

	setupStaticFiles(e)

	setupErrorHandler(e)

	startServer(e)
}

func setupCORS(e *echo.Echo) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusOK)
			}

			return next(c)
		}
	})
}

func setupRoutes(e *echo.Echo, container *container.Container, jwt string) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":    "ok",
			"message":   "MicroShopik API is running",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
		})
	})

	setupAuthRoutes(e, container)

	setupCategoryRoutes(e, container, jwt)

	setupProductRoutes(e, container, jwt)

	setupOrderRoutes(e, container, jwt)

	setupConversationRoutes(e, container, jwt)

	setupRoleRoutes(e, container, jwt)

	setupAdminRoutes(e, container, jwt)

	setupSellerRoutes(e, jwt)
}

func setupAuthRoutes(e *echo.Echo, container *container.Container) {
	auth := e.Group("/auth")
	auth.POST("/register", container.UserController.Register)
	auth.POST("/login", container.UserController.Login)
}

func setupCategoryRoutes(e *echo.Echo, container *container.Container, jwt string) {
	cats := e.Group("/categories")
	cats.GET("", container.CategoryController.GetAll)
	cats.GET("/:id", container.CategoryController.GetById)

	catsAdmin := e.Group("/categories")
	catsAdmin.Use(middleware.JWTMiddleware(jwt))
	catsAdmin.Use(middleware.RequireRole("admin"))
	catsAdmin.POST("", container.CategoryController.Create)
	catsAdmin.PUT("/:id", container.CategoryController.Update)
	catsAdmin.DELETE("/:id", container.CategoryController.Delete)
}

func setupProductRoutes(e *echo.Echo, container *container.Container, jwt string) {
	products := e.Group("/products")
	products.GET("", container.ProductController.Find)
	products.GET("/count", container.ProductController.Count)
	products.GET("/:id", container.ProductController.GetById)
	products.GET("/:id/available", container.ProductController.IsAvailable)

	productsAuth := e.Group("/products")
	productsAuth.Use(middleware.JWTMiddleware(jwt))
	productsAuth.Use(middleware.RequireAnyRole("seller", "admin"))
	productsAuth.POST("", container.ProductController.Create)
	productsAuth.PUT("/:id", container.ProductController.Update)
	productsAuth.DELETE("/:id", container.ProductController.Delete)
}

func setupOrderRoutes(e *echo.Echo, container *container.Container, jwt string) {
	orders := e.Group("/orders")
	orders.Use(middleware.JWTMiddleware(jwt))
	orders.POST("", container.OrderController.Create)
	orders.GET("", container.OrderController.GetMyOrders)
	orders.GET("/seller", container.OrderController.GetMyOrdersAsSeller)
	orders.GET("/:id", container.OrderController.GetByID)
	orders.PUT("/:id/status", container.OrderController.UpdateStatus)
	orders.POST("/:id/process", container.OrderController.ProcessOrder)
	orders.POST("/:id/cancel/:customerID", container.OrderController.CancelOrder)
	orders.POST("/:id/confirm", container.OrderController.ConfirmOrder)
}

func setupConversationRoutes(e *echo.Echo, container *container.Container, jwt string) {
	conversations := e.Group("/conversations")
	conversations.Use(middleware.JWTMiddleware(jwt))
	conversations.POST("", container.ConversationController.Create)
	conversations.GET("/:id", container.ConversationController.GetByID)
	conversations.PUT("/:id", container.ConversationController.Update)
	conversations.DELETE("/:id", container.ConversationController.Delete)
	conversations.POST("/:id/participants/:userID", container.ConversationController.AddParticipant)
	conversations.DELETE("/:id/participants/:userID", container.ConversationController.RemoveParticipant)

	users := e.Group("/users")
	users.Use(middleware.JWTMiddleware(jwt))
	users.GET("/:userID/conversations", container.ConversationController.GetByUserID)

	messages := e.Group("/conversations/:conversationID/messages")
	messages.Use(middleware.JWTMiddleware(jwt))
	messages.GET("", container.MessageController.GetByConversationID)
	messages.POST("", container.MessageController.Create)
	messages.GET("/system", container.MessageController.GetSystemMessages)
	messages.POST("/system", container.MessageController.SendSystemMessage)
}

func setupRoleRoutes(e *echo.Echo, container *container.Container, jwt string) {
	roles := e.Group("/roles")
	roles.Use(middleware.JWTMiddleware(jwt))
	roles.Use(middleware.RequireRole("admin"))
	roles.GET("", container.RoleController.GetAll)
	roles.GET("/:name", container.RoleController.GetByName)
	roles.POST("", container.RoleController.Create)
	roles.PUT("/:id", container.RoleController.Update)
	roles.DELETE("/:id", container.RoleController.Delete)

	userManagement := e.Group("/users")
	userManagement.Use(middleware.JWTMiddleware(jwt))
	userManagement.Use(middleware.RequireRole("admin"))
	userManagement.GET("/:user_id/roles", container.RoleController.GetUserRoles)
	userManagement.POST("/:user_id/roles/:role_name", container.RoleController.AssignRoleToUser)
	userManagement.DELETE("/:user_id/roles/:role_name", container.RoleController.RemoveRoleFromUser)
	userManagement.PUT("/:id", container.UserController.AdminUpdate)
}

func setupAdminRoutes(e *echo.Echo, container *container.Container, jwt string) {
	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.JWTMiddleware(jwt))
	adminGroup.Use(middleware.RequireRole("admin"))

	adminGroup.GET("/dashboard", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Admin dashboard - you can manage all users and products",
		})
	})

	adminGroup.GET("/users", func(c echo.Context) error {
		users, err := container.UserRepository.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get users"})
		}
		return c.JSON(200, users)
	})
	adminGroup.GET("/users/:id", container.UserController.AdminGet)
	adminGroup.PUT("/users/:id", container.UserController.AdminUpdate)
	adminGroup.GET("/users/:user_id/roles", container.RoleController.GetUserRoles)
	adminGroup.POST("/users/:user_id/roles/:role_name", container.RoleController.AssignRoleToUser)
	adminGroup.DELETE("/users/:user_id/roles/:role_name", container.RoleController.RemoveRoleFromUser)

	adminGroup.GET("/products", func(c echo.Context) error {
		products, err := container.ProductRepository.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get products"})
		}
		return c.JSON(200, products)
	})

	adminGroup.GET("/orders", func(c echo.Context) error {
		orders, err := container.OrderRepository.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get orders"})
		}
		return c.JSON(200, orders)
	})

	adminGroup.GET("/stats", func(c echo.Context) error {
		users, _ := container.UserRepository.GetAll()
		products, _ := container.ProductRepository.GetAll()
		orders, _ := container.OrderRepository.GetAll()

		stats := map[string]interface{}{
			"total_users":     len(users),
			"total_products":  len(products),
			"total_orders":    len(orders),
			"database_status": "Online",
			"api_status":      "Running",
			"system_status":   "operational",
		}
		return c.JSON(200, stats)
	})
}

func setupSellerRoutes(e *echo.Echo, jwt string) {
	sellerGroup := e.Group("/seller")
	sellerGroup.Use(middleware.JWTMiddleware(jwt))
	sellerGroup.Use(middleware.RequireRole("seller"))

	sellerGroup.GET("/products", func(c echo.Context) error {
		userID := c.Get("user_id").(int)
		return c.JSON(200, map[string]interface{}{
			"message": "Seller dashboard - you can manage your products here",
			"user_id": userID,
		})
	})
}

func setupStaticFiles(e *echo.Echo) {
	e.Static("/assets", "frontend/dist/assets")
	e.File("/", "frontend/dist/index.html")
}

func setupErrorHandler(e *echo.Echo) {
	originalHTTPErrorHandler := e.HTTPErrorHandler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var he *echo.HTTPError
		if errors.As(err, &he) && he.Code == http.StatusNotFound {
			req := c.Request()
			if req.Method == http.MethodGet && strings.Contains(req.Header.Get("Accept"), "text/html") {
				_ = c.File("frontend/dist/index.html")
				return
			}
		}
		if originalHTTPErrorHandler != nil {
			originalHTTPErrorHandler(err, c)
			return
		}
		_ = c.JSON(http.StatusInternalServerError, map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
	}
}

func startServer(e *echo.Echo) {
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

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

	err = database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	e := echo.New()

	userRepo := repositories.NewUserRepository(database.GetDB())
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authCtrl := controllers.NewAuthController(authService)

	roleRepo := repositories.NewRoleRepository(database.GetDB())
	roleService := services.NewRoleService(roleRepo, userRepo)
	roleCtrl := controllers.NewRoleController(roleService)

	productRepo := repositories.NewProductRepository(database.GetDB())
	productService := services.NewProductService(productRepo)
	productCtrl := controllers.NewProductController(productService)

	catRepo := repositories.NewCategoryRepository(database.GetDB())
	catService := services.NewCategoryService(catRepo)
	catCtrl := controllers.NewCategoryController(catService)

	convRepo := repositories.NewConversationRepository(database.GetDB())

	participantRepo := repositories.NewParticipantRepository(database.GetDB())
	participantService := services.NewParticipantService(participantRepo, convRepo, userRepo)
	convService := services.NewConversationService(convRepo, participantRepo, userRepo)

	convCtrl := controllers.NewConversationController(convService)
	participantCtrl := controllers.NewParticipantController(participantService)

	pItemRepo := repositories.NewProductItemRepository(database.GetDB())
	pItemService := services.NewProductItemService(pItemRepo, productRepo)
	pItemCtrl := controllers.NewProductItemController(pItemService)

	orderRepo := repositories.NewOrderRepository(database.GetDB())
	orderService := services.NewOrderService(orderRepo, pItemRepo, productRepo, userRepo, convService)
	orderCtrl := controllers.NewOrderController(orderService)

	messageRepo := repositories.NewMessageRepository(database.GetDB())
	messageService := services.NewMessageService(messageRepo, convRepo, participantRepo, orderRepo)
	messageCtrl := controllers.NewMessageController(messageService)

	orders := e.Group("/orders")
	orders.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	orders.POST("", orderCtrl.Create)
	orders.GET("", func(c echo.Context) error {
		userID := c.Get("user_id").(int)
		return orderCtrl.GetByCustomerID(c)
	})
	orders.GET("/:id", orderCtrl.GetByID)
	orders.PUT("/:id", orderCtrl.Update)
	orders.DELETE("/:id", orderCtrl.Delete)
	orders.PUT("/:id/status", orderCtrl.UpdateStatus)
	orders.POST("/:id/process", orderCtrl.ProcessOrder)
	orders.POST("/:id/cancel/:customerID", orderCtrl.CancelOrder)

	conversations := e.Group("/conversations")
	conversations.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	conversations.POST("", convCtrl.Create)
	conversations.GET("/:id", convCtrl.GetByID)
	conversations.PUT("/:id", convCtrl.Update)
	conversations.DELETE("/:id", convCtrl.Delete)
	conversations.POST("/:id/participants/:userID", convCtrl.AddParticipant)
	conversations.DELETE("/:id/participants/:userID", convCtrl.RemoveParticipant)

	users := e.Group("/users")
	users.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	users.GET("/:userID/conversations", convCtrl.GetByUserID)

	messages := e.Group("/conversations/:conversationID/messages")
	messages.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	messages.GET("", messageCtrl.GetByConversationID)
	messages.POST("", messageCtrl.Create)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok", "message": "MicroShopik API is running"})
	})

	auth := e.Group("/auth")
	auth.POST("/register", authCtrl.Register)
	auth.POST("/login", authCtrl.Login)

	cats := e.Group("/categories")
	cats.GET("", catCtrl.GetAll)
	cats.GET("/:id", catCtrl.GetById)

	catsAdmin := e.Group("/categories")
	catsAdmin.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	catsAdmin.Use(middleware.RequireRole("admin"))
	catsAdmin.POST("", catCtrl.Create)
	catsAdmin.PUT("/:id", catCtrl.Update)
	catsAdmin.DELETE("/:id", catCtrl.Delete)

	products := e.Group("/products")
	products.GET("", productCtrl.Find)
	products.GET("/count", productCtrl.Count)
	products.GET("/:id", productCtrl.GetById)
	products.GET("/:id/available", productCtrl.IsAvailable)

	productsAuth := e.Group("/products")
	productsAuth.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	productsAuth.Use(middleware.RequireAnyRole("seller", "admin"))
	productsAuth.POST("", productCtrl.Create)
	productsAuth.PUT("/:id", productCtrl.Update)
	productsAuth.DELETE("/:id", productCtrl.Delete)

	roles := e.Group("/roles")
	roles.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	roles.Use(middleware.RequireRole("admin"))
	roles.GET("", roleCtrl.GetAll)
	roles.GET("/:name", roleCtrl.GetByName)
	roles.POST("", roleCtrl.Create)

	users := e.Group("/users")
	users.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	users.Use(middleware.RequireRole("admin"))
	users.GET("/:user_id/roles", roleCtrl.GetUserRoles)
	users.POST("/:user_id/roles/:role_name", roleCtrl.AssignRoleToUser)

	sellerGroup := e.Group("/seller")
	sellerGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	sellerGroup.Use(middleware.RequireRole("seller"))
	sellerGroup.GET("/products", func(c echo.Context) error {
		userID := c.Get("user_id").(int)
		return c.JSON(200, map[string]interface{}{
			"message": "Seller dashboard - you can manage your products here",
			"user_id": userID,
		})
	})

	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	adminGroup.Use(middleware.RequireRole("admin"))
	adminGroup.GET("/dashboard", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Admin dashboard - you can manage all users and products",
		})
	})
	adminGroup.GET("/users", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "List all users (admin only)",
		})
	})

	moderatorGroup := e.Group("/moderator")
	moderatorGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	moderatorGroup.Use(middleware.RequireAnyRole("admin", "moderator"))
	moderatorGroup.GET("/reports", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "View reports (admin or moderator)",
		})
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

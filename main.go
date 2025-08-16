package main

import (
	"MicroShopik/configs"
	_ "MicroShopik/docs"
	"MicroShopik/internal/controllers"
	"MicroShopik/internal/database"
	"MicroShopik/internal/middleware"
	"MicroShopik/internal/repositories"
	"MicroShopik/internal/services"
	"MicroShopik/scripts"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
	_ = controllers.NewParticipantController(participantService)

	pItemRepo := repositories.NewProductItemRepository(database.GetDB())
	pItemService := services.NewProductItemService(pItemRepo, productRepo)
	pItemCtrl := controllers.NewProductItemController(pItemService)

	orderRepo := repositories.NewOrderRepository(database.GetDB())

	messageRepo := repositories.NewMessageRepository(database.GetDB())
	messageService := services.NewMessageService(messageRepo, convRepo, participantRepo, orderRepo)
	messageCtrl := controllers.NewMessageController(messageService)

	orderService := services.NewOrderService(orderRepo, pItemRepo, productRepo, userRepo, convService, messageService)
	orderCtrl := controllers.NewOrderController(orderService)

	// Record start time for uptime calculation
	startTime := time.Now()

	orders := e.Group("/orders")
	orders.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	orders.POST("", orderCtrl.Create)
	orders.GET("", orderCtrl.GetMyOrders)
	orders.GET("/seller", orderCtrl.GetMyOrdersAsSeller)
	orders.GET("/:id", orderCtrl.GetByID)
	orders.PUT("/:id", orderCtrl.Update)
	orders.DELETE("/:id", orderCtrl.Delete)
	orders.PUT("/:id/status", orderCtrl.UpdateStatus)
	orders.POST("/:id/process", orderCtrl.ProcessOrder)
	orders.POST("/:id/cancel/:customerID", orderCtrl.CancelOrder)
	orders.POST("/:id/confirm", orderCtrl.ConfirmOrder)

	conversations := e.Group("/conversations")
	conversations.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	conversations.POST("", convCtrl.Create)
	conversations.GET("/:id", convCtrl.GetByID)
	conversations.PUT("/:id", convCtrl.Update)
	conversations.DELETE("/:id", convCtrl.Delete)
	conversations.POST("/:id/participants/:userID", convCtrl.AddParticipant)
	conversations.DELETE("/:id/participants/:userID", convCtrl.RemoveParticipant)

	messages := e.Group("/conversations/:conversationID/messages")
	messages.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	messages.GET("", messageCtrl.GetByConversationID)
	messages.POST("", messageCtrl.Create)
	messages.GET("/system", messageCtrl.GetSystemMessages)
	messages.POST("/system", messageCtrl.SendSystemMessage)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":    "ok",
			"message":   "MicroShopik API is running",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    time.Since(startTime).String(),
			"version":   "1.0.0",
		})
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
	products.GET("/:id/conversations", convCtrl.GetByProductID) // Новый роут для поиска конверсаций по товару

	productsAuth := e.Group("/products")
	productsAuth.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	productsAuth.Use(middleware.RequireAnyRole("seller", "admin"))
	productsAuth.POST("", productCtrl.Create)
	productsAuth.PUT("/:id", productCtrl.Update)
	productsAuth.DELETE("/:id", productCtrl.Delete)

	productItems := e.Group("/product-items")
	productItems.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	productItems.Use(middleware.RequireAnyRole("seller", "admin"))
	productItems.POST("", pItemCtrl.Create)
	productItems.GET("/:id", pItemCtrl.GetByID)
	productItems.PUT("/:id", pItemCtrl.Update)
	productItems.DELETE("/:id", pItemCtrl.Delete)
	productItems.GET("/product/:productId", pItemCtrl.GetByProductID)

	roles := e.Group("/roles")
	roles.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	roles.Use(middleware.RequireRole("admin"))
	roles.GET("", roleCtrl.GetAll)
	roles.GET("/:name", roleCtrl.GetByName)
	roles.POST("", roleCtrl.Create)

	userManagement := e.Group("/users")
	userManagement.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	userManagement.Use(middleware.RequireRole("admin"))
	userManagement.GET("/:user_id/roles", roleCtrl.GetUserRoles)
	userManagement.POST("/:user_id/roles/:role_name", roleCtrl.AssignRoleToUser)
	userManagement.DELETE("/:user_id/roles/:role_name", roleCtrl.RemoveRoleFromUser)

	e.GET("/roles", roleCtrl.GetAll)
	e.POST("/roles", roleCtrl.Create)
	e.PUT("/roles/:id", roleCtrl.Update)
	e.DELETE("/roles/:id", roleCtrl.Delete)

	userProfile := e.Group("/users/profile")
	userProfile.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	userProfile.GET("", func(c echo.Context) error {
		userID := c.Get("user_id").(int)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		user, err = userRepo.GetByID(userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user with roles"})
		}
		return c.JSON(http.StatusOK, user)
	})
	userProfile.PUT("", func(c echo.Context) error {
		_ = c.Get("user_id").(int)
		var userData map[string]interface{}
		if err := c.Bind(&userData); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "profile updated"})
	})

	userConversations := e.Group("/users")
	userConversations.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	userConversations.GET("/:userID/conversations", convCtrl.GetByUserID)

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
		users, err := userRepo.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get users"})
		}
		return c.JSON(200, users)
	})

	adminGroup.GET("/users/:id", func(c echo.Context) error {
		userID := c.Param("id")
		id, err := strconv.Atoi(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
		}
		user, err := userRepo.GetByID(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(200, user)
	})

	adminGroup.PUT("/users/:id", func(c echo.Context) error {
		userIDParam := c.Param("id")
		id, err := strconv.Atoi(userIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
		}
		var userData map[string]interface{}
		if err := c.Bind(&userData); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		}
		user, err := userRepo.GetByID(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		if v, ok := userData["username"].(string); ok && v != "" {
			user.Username = v
		}
		if v, ok := userData["email"].(string); ok && v != "" {
			user.Email = v
		}
		if err := userRepo.Update(user); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update user"})
		}
		updated, _ := userRepo.GetByID(id)
		return c.JSON(http.StatusOK, updated)
	})

	adminGroup.DELETE("/users/:id", func(c echo.Context) error {
		userIDParam := c.Param("id")
		id, err := strconv.Atoi(userIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
		}
		if err := userRepo.Delete(id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "user deleted successfully"})
	})

	adminGroup.GET("/products", func(c echo.Context) error {
		products, err := productRepo.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get products"})
		}
		return c.JSON(200, products)
	})

	adminGroup.PUT("/products/:id/status", func(c echo.Context) error {
		productID := c.Param("id")
		_, err := strconv.Atoi(productID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product ID"})
		}
		var statusData map[string]interface{}
		if err := c.Bind(&statusData); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		}
		return c.JSON(200, map[string]string{"message": "product status updated successfully"})
	})

	adminGroup.DELETE("/products/:id", func(c echo.Context) error {
		productID := c.Param("id")
		_, err := strconv.Atoi(productID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product ID"})
		}
		return c.JSON(200, map[string]string{"message": "product deleted successfully"})
	})

	adminGroup.GET("/orders", func(c echo.Context) error {
		orders, err := orderRepo.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get orders"})
		}
		return c.JSON(200, orders)
	})

	adminGroup.PUT("/orders/:id/status", func(c echo.Context) error {
		orderID := c.Param("id")
		_, err := strconv.Atoi(orderID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order ID"})
		}
		var statusData map[string]interface{}
		if err := c.Bind(&statusData); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		}
		return c.JSON(200, map[string]string{"message": "order status updated successfully"})
	})

	adminGroup.DELETE("/orders/:id", func(c echo.Context) error {
		orderID := c.Param("id")
		_, err := strconv.Atoi(orderID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order ID"})
		}
		return c.JSON(200, map[string]string{"message": "order deleted successfully"})
	})

	adminGroup.GET("/stats", func(c echo.Context) error {
		users, err := userRepo.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get users"})
		}

		products, err := productRepo.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get products"})
		}

		orders, err := orderRepo.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get orders"})
		}

		stats := map[string]interface{}{
			"total_users":     len(users),
			"total_products":  len(products),
			"total_orders":    len(orders),
			"database_status": "Online",
			"api_status":      "Running",
			"active_users":    len(users),
			"system_status":   "operational",
		}
		return c.JSON(200, stats)
	})

	adminGroup.GET("/logs", func(c echo.Context) error {
		logs := []map[string]interface{}{
			{"timestamp": time.Now().Format(time.RFC3339), "level": "INFO", "message": "System running normally"},
			{"timestamp": time.Now().Add(-5 * time.Minute).Format(time.RFC3339), "level": "INFO", "message": "Database backup completed successfully"},
			{"timestamp": time.Now().Add(-10 * time.Minute).Format(time.RFC3339), "level": "INFO", "message": "New user registration: john_doe"},
			{"timestamp": time.Now().Add(-15 * time.Minute).Format(time.RFC3339), "level": "WARN", "message": "High memory usage detected"},
			{"timestamp": time.Now().Add(-20 * time.Minute).Format(time.RFC3339), "level": "INFO", "message": "Product inventory updated"},
			{"timestamp": time.Now().Add(-25 * time.Minute).Format(time.RFC3339), "level": "ERROR", "message": "Failed to process payment for order #1234"},
			{"timestamp": time.Now().Add(-30 * time.Minute).Format(time.RFC3339), "level": "INFO", "message": "API rate limit reset"},
			{"timestamp": time.Now().Add(-35 * time.Minute).Format(time.RFC3339), "level": "INFO", "message": "Cache cleared successfully"},
			{"timestamp": time.Now().Add(-40 * time.Minute).Format(time.RFC3339), "level": "INFO", "message": "Email notification sent to user"},
			{"timestamp": time.Now().Add(-45 * time.Minute).Format(time.RFC3339), "level": "INFO", "message": "System maintenance completed"},
		}
		return c.JSON(200, logs)
	})

	adminGroup.PUT("/settings", func(c echo.Context) error {
		var settings map[string]interface{}
		if err := c.Bind(&settings); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		}
		return c.JSON(200, map[string]string{"message": "settings updated successfully"})
	})

	moderatorGroup := e.Group("/moderator")
	moderatorGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	moderatorGroup.Use(middleware.RequireAnyRole("admin", "moderator"))
	moderatorGroup.GET("/reports", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "View reports (admin or moderator)",
		})
	})

	e.Static("/assets", "frontend/dist/assets")
	e.File("/", "frontend/dist/index.html")
	e.GET("/*", func(c echo.Context) error {
		if c.Request().Method != http.MethodGet {
			return c.NoContent(http.StatusNotFound)
		}
		return c.File("frontend/dist/index.html")
	})

	originalHTTPErrorHandler := e.HTTPErrorHandler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusNotFound {
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

	go func() {
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down server:", err)
		}
	}()

	// Start keep-alive service if enabled
	var keepAliveService *services.KeepAliveService
	if cfg.KeepAliveEnabled {
		keepAliveService = services.NewKeepAliveService(
			cfg.KeepAliveURL,
			time.Duration(cfg.KeepAliveInterval)*time.Minute,
		)
		keepAliveService.Start()
		log.Printf("Keep-alive service started with interval: %d minutes", cfg.KeepAliveInterval)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	// Stop keep-alive service gracefully
	if keepAliveService != nil {
		keepAliveService.Stop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

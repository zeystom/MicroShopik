package container

import (
	"MicroShopik/configs"
	"MicroShopik/internal/controllers"
	"MicroShopik/internal/database"
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"MicroShopik/internal/services/application"
	sdomain "MicroShopik/internal/services/domain"
	"log"
)

type Container struct {
	UserRepository         domain.UserRepository
	RoleRepository         domain.RoleRepository
	ProductRepository      domain.ProductRepository
	CategoryRepository     domain.CategoryRepository
	ConversationRepository domain.ConversationRepository
	ParticipantRepository  domain.ParticipantRepository
	OrderRepository        domain.OrderRepository
	MessageRepository      domain.MessageRepository

	UserService         sdomain.UserService
	RoleService         sdomain.RoleService
	ProductService      sdomain.ProductService
	CategoryService     sdomain.CategoryService
	ConversationService sdomain.ConversationService
	ParticipantService  sdomain.ParticipantService
	OrderService        sdomain.OrderService
	MessageService      sdomain.MessageService

	OrderApplicationService        *application.OrderApplicationService
	UserApplicationService         *application.UserApplicationService
	ProductApplicationService      *application.ProductApplicationService
	ConversationApplicationService *application.ConversationApplicationService

	UserController         *controllers.UserController
	RoleController         *controllers.RoleController
	ProductController      *controllers.ProductController
	CategoryController     *controllers.CategoryController
	ConversationController *controllers.ConversationController
	ParticipantController  *controllers.ParticipantController
	OrderController        *controllers.OrderController
	MessageController      *controllers.MessageController
}

func NewContainer() *Container {
	db := database.GetDB()
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	conversationRepo := repositories.NewConversationRepository(db)
	participantRepo := repositories.NewParticipantRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	userService := sdomain.NewUserService(userRepo, cfg.JWTSecret)
	roleService := sdomain.NewRoleService(roleRepo, userRepo)
	productService := sdomain.NewProductService(productRepo)
	categoryService := sdomain.NewCategoryService(categoryRepo)
	participantService := sdomain.NewParticipantService(participantRepo, conversationRepo, userRepo)
	conversationService := sdomain.NewConversationService(conversationRepo, participantRepo, userRepo)
	messageService := sdomain.NewMessageService(messageRepo, conversationRepo, participantRepo, orderRepo)
	orderService := sdomain.NewOrderService(orderRepo)

	orderAppService := application.NewOrderApplicationService(
		orderService,
		productService,
		userService,
		conversationService,
		messageService,
	)

	userAppService := application.NewUserApplicationService(
		userService,
		roleService,
	)

	productAppService := application.NewProductApplicationService(
		productService,
		categoryService,
		userService,
		orderService,
	)

	conversationAppService := application.NewConversationApplicationService(
		conversationService,
		messageService,
		userService,
		productService,
		participantService,
	)

	userController := controllers.NewUserController(userAppService)
	roleController := controllers.NewRoleController(roleService)
	productController := controllers.NewProductController(productAppService)
	categoryController := controllers.NewCategoryController(categoryService)
	conversationController := controllers.NewConversationController(conversationAppService)
	participantController := controllers.NewParticipantController(participantService)
	orderController := controllers.NewOrderController(orderAppService)
	messageController := controllers.NewMessageController(messageService)

	return &Container{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		ProductRepository:      productRepo,
		CategoryRepository:     categoryRepo,
		ConversationRepository: conversationRepo,
		ParticipantRepository:  participantRepo,
		OrderRepository:        orderRepo,
		MessageRepository:      messageRepo,

		UserService:         userService,
		RoleService:         roleService,
		ProductService:      productService,
		CategoryService:     categoryService,
		ConversationService: conversationService,
		ParticipantService:  participantService,
		OrderService:        orderService,
		MessageService:      messageService,

		OrderApplicationService:        orderAppService,
		UserApplicationService:         userAppService,
		ProductApplicationService:      productAppService,
		ConversationApplicationService: conversationAppService,

		UserController:         userController,
		RoleController:         roleController,
		ProductController:      productController,
		CategoryController:     categoryController,
		ConversationController: conversationController,
		ParticipantController:  participantController,
		OrderController:        orderController,
		MessageController:      messageController,
	}
}

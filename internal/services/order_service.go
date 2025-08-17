package services

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"errors"
)

type OrderService interface {
	Create(order *domain.Order) error
	GetByID(id int) (*domain.Order, error)
	GetByCustomerID(customerID int) ([]*domain.Order, error)
	GetBySellerID(sellerID int) ([]*domain.Order, error)
	GetByStatus(status string) ([]*domain.Order, error)
	Update(order *domain.Order) error
	Delete(id int) error
	UpdateStatus(id int, status string) error
	GetByProductID(productID int) ([]*domain.Order, error)
	ProcessOrder(orderID int) error
	CancelOrder(orderID int, customerID int) error
	ConfirmOrder(orderID int, customerID int) error
}

type orderService struct {
	orderRepo           repositories.OrderRepository
	productItemRepo     repositories.ProductItemRepository
	productRepo         repositories.ProductCRUDRepository
	userRepo            repositories.UserRepository
	conversationService ConversationService
	messageService      MessageService
}

func NewOrderService(oRepo repositories.OrderRepository, piRepo repositories.ProductItemRepository, pRepo repositories.ProductCRUDRepository, uRepo repositories.UserRepository, convService ConversationService, msgService MessageService) OrderService {
	return &orderService{
		orderRepo:           oRepo,
		productItemRepo:     piRepo,
		productRepo:         pRepo,
		userRepo:            uRepo,
		conversationService: convService,
		messageService:      msgService,
	}
}

func (s *orderService) Create(order *domain.Order) error {
	if order.CustomerID != nil {
		_, err := s.userRepo.GetByID(*order.CustomerID)
		if err != nil {
			return errors.New("customer not found")
		}
	}

	// Validate product and prevent self-purchase
	if order.ProductID != nil {
		product, err := s.productRepo.GetById(*order.ProductID)
		if err != nil {
			return errors.New("product not found")
		}
		if order.CustomerID != nil && product.SellerID == *order.CustomerID {
			return errors.New("cannot purchase your own product")
		}
	}

	if order.ProductItemID != nil {
		productItem, err := s.productItemRepo.GetByID(*order.ProductItemID)
		if err != nil {
			return errors.New("product item not found")
		}
		if productItem.IsUsed {
			return errors.New("product item is already used")
		}
		// Extra safety: prevent self-purchase via product item as well
		if order.CustomerID != nil && productItem.Product.ID != 0 && productItem.Product.SellerID == *order.CustomerID {
			return errors.New("cannot purchase your own product")
		}
	}

	return s.orderRepo.Create(order)
}

func (s *orderService) GetByID(id int) (*domain.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *orderService) GetByCustomerID(customerID int) ([]*domain.Order, error) {
	return s.orderRepo.GetByCustomerID(customerID)
}

func (s *orderService) GetBySellerID(sellerID int) ([]*domain.Order, error) {
	return s.orderRepo.GetBySellerID(sellerID)
}

func (s *orderService) GetByStatus(status string) ([]*domain.Order, error) {
	return s.orderRepo.GetByStatus(status)
}

func (s *orderService) Update(order *domain.Order) error {
	return s.orderRepo.Update(order)
}

func (s *orderService) Delete(id int) error {
	return s.orderRepo.Delete(id)
}

func (s *orderService) UpdateStatus(id int, status string) error {
	validStatuses := []string{"pending", "completed", "refunded", "cancelled"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid status")
	}

	return s.orderRepo.UpdateStatus(id, status)
}

func (s *orderService) GetByProductID(productID int) ([]*domain.Order, error) {
	return s.orderRepo.GetByProductID(productID)
}

func (s *orderService) ProcessOrder(orderID int) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != "pending" {
		return errors.New("order is not in pending status")
	}

	tx := s.orderRepo.BeginTx()
	defer tx.Rollback()

	if order.ProductID != nil {
		isAvailable, err := s.productRepo.CheckAvailabilityAndIncrementSoldCount(*order.ProductID, 1)
		if err != nil {
			return err
		}
		if !isAvailable {
			return errors.New("product is not available or has reached maximum sales limit")
		}
	}

	if order.ProductItemID != nil {
		if err := s.productItemRepo.MarkAsUsedTx(tx, *order.ProductItemID); err != nil {
			return err
		}
	}

	if err := s.orderRepo.UpdateStatusTx(tx, orderID, "completed"); err != nil {
		return err
	}

	if order.CustomerID != nil {
		product, err := s.productRepo.GetById(*order.ProductID)
		if err != nil {
			return err
		}

		conversation := &domain.Conversation{}
		participantIDs := []int{*order.CustomerID, product.SellerID}

		if err := s.conversationService.Create(conversation, participantIDs); err != nil {
			return err
		}

		if err := s.messageService.SendSystemMessage(conversation.ID, "Order has been processed successfully. You can now discuss delivery details.", &orderID); err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (s *orderService) CancelOrder(orderID int, customerID int) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.CustomerID == nil || *order.CustomerID != customerID {
		return errors.New("unauthorized to cancel this order")
	}

	if order.Status != "pending" {
		return errors.New("order cannot be cancelled")
	}

	return s.orderRepo.UpdateStatus(orderID, "cancelled")
}

func (s *orderService) ConfirmOrder(orderID int, customerID int) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.CustomerID == nil || *order.CustomerID != customerID {
		return errors.New("unauthorized to confirm this order")
	}

	if order.Status != "pending" {
		return errors.New("order cannot be confirmed")
	}

	if err := s.orderRepo.UpdateStatus(orderID, "confirmed"); err != nil {
		return err
	}

	if order.ProductID != nil {
		isAvailable, err := s.productRepo.CheckAvailabilityAndIncrementSoldCount(*order.ProductID, 1)
		if err != nil {
			return err
		}
		if !isAvailable {
			if err := s.orderRepo.UpdateStatus(orderID, "cancelled"); err != nil {
				return err
			}
			return errors.New("product is no longer available, order has been cancelled")
		}
	}

	return nil
}

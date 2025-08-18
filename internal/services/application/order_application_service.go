package application

import (
	"MicroShopik/internal/domain"
	domain2 "MicroShopik/internal/services/domain"
	"errors"
)

type OrderApplicationService struct {
	orderService        domain2.OrderService
	productService      domain2.ProductService
	userService         domain2.UserService
	conversationService domain2.ConversationService
	messageService      domain2.MessageService
}

func NewOrderApplicationService(
	orderService domain2.OrderService,
	productService domain2.ProductService,
	userService domain2.UserService,
	conversationService domain2.ConversationService,
	messageService domain2.MessageService,
) *OrderApplicationService {
	return &OrderApplicationService{
		orderService:        orderService,
		productService:      productService,
		userService:         userService,
		conversationService: conversationService,
		messageService:      messageService,
	}
}

func (s *OrderApplicationService) CreateOrder(order *domain.Order) error {
	if order.CustomerID != nil {
		if err := s.userService.ValidateUserExists(*order.CustomerID); err != nil {
			return errors.New("customer not found")
		}
	}

	if order.ProductID != nil {
		if err := s.productService.ValidateProductForOrder(*order.ProductID, order.CustomerID); err != nil {
			return err
		}
	}

	return s.orderService.Create(order)
}

func (s *OrderApplicationService) ProcessOrder(orderID int) error {
	order, err := s.orderService.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != "pending" {
		return errors.New("order is not in pending status")
	}

	if order.ProductID != nil {
		if err := s.productService.ReserveProduct(*order.ProductID); err != nil {
			return err
		}
	}

	if err := s.orderService.UpdateStatus(orderID, "completed"); err != nil {
		return err
	}

	if order.CustomerID != nil && order.ProductID != nil {
		if err := s.createOrderConversation(order); err != nil {
			return err
		}
	}

	return nil
}

func (s *OrderApplicationService) CancelOrder(orderID int, customerID int) error {
	order, err := s.orderService.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.CustomerID == nil || *order.CustomerID != customerID {
		return errors.New("unauthorized to cancel this order")
	}

	if order.Status != "pending" {
		return errors.New("order cannot be cancelled")
	}

	if order.ProductID != nil {
		if err := s.productService.ReleaseProduct(*order.ProductID); err != nil {
			return err
		}
	}

	return s.orderService.UpdateStatus(orderID, "cancelled")
}

func (s *OrderApplicationService) ConfirmOrder(orderID int, customerID int) error {
	order, err := s.orderService.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.CustomerID == nil || *order.CustomerID != customerID {
		return errors.New("unauthorized to confirm this order")
	}

	if order.Status != "pending" {
		return errors.New("order cannot be confirmed")
	}

	if order.ProductID != nil {
		if err := s.productService.ReserveProduct(*order.ProductID); err != nil {
			return err
		}
	}

	return s.orderService.UpdateStatus(orderID, "confirmed")
}

func (s *OrderApplicationService) GetMyOrders(userID int) ([]*domain.Order, error) {
	return s.orderService.GetByCustomerID(userID)
}

func (s *OrderApplicationService) GetMyOrdersAsSeller(sellerID int) ([]*domain.Order, error) {
	return s.orderService.GetBySellerID(sellerID)
}

func (s *OrderApplicationService) GetOrderByID(orderID int) (*domain.Order, error) {
	return s.orderService.GetByID(orderID)
}

func (s *OrderApplicationService) UpdateOrderStatus(orderID int, status string) error {
	return s.orderService.UpdateStatus(orderID, status)
}

func (s *OrderApplicationService) createOrderConversation(order *domain.Order) error {
	if order.ProductID == nil || order.CustomerID == nil {
		return nil
	}

	product, err := s.productService.GetById(*order.ProductID)
	if err != nil {
		return err
	}

	conversation := &domain.Conversation{}
	participantIDs := []int{*order.CustomerID, product.SellerID}

	if err := s.conversationService.Create(conversation, participantIDs); err != nil {
		return err
	}

	if err := s.messageService.SendSystemMessage(
		conversation.ID,
		"Order has been processed successfully. You can now discuss delivery details.",
		&order.ID,
	); err != nil {
		return err
	}

	return nil
}

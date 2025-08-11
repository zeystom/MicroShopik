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
	GetByStatus(status string) ([]*domain.Order, error)
	Update(order *domain.Order) error
	Delete(id int) error
	UpdateStatus(id int, status string) error
	GetByProductID(productID int) ([]*domain.Order, error)
	ProcessOrder(orderID int) error
	CancelOrder(orderID int, customerID int) error
}

type orderService struct {
	orderRepo       repositories.OrderRepository
	productItemRepo repositories.ProductItemRepository
	productRepo     repositories.ProductCRUDRepository
	userRepo        repositories.UserRepository
}

func NewOrderService(oRepo repositories.OrderRepository, piRepo repositories.ProductItemRepository, pRepo repositories.ProductCRUDRepository, uRepo repositories.UserRepository) OrderService {
	return &orderService{
		orderRepo:       oRepo,
		productItemRepo: piRepo,
		productRepo:     pRepo,
		userRepo:        uRepo,
	}
}

func (s *orderService) Create(order *domain.Order) error {
	// Validate customer exists
	if order.CustomerID != nil {
		_, err := s.userRepo.GetByEmail("") // We need a GetByID method, using placeholder for now
		if err != nil {
			return errors.New("customer not found")
		}
	}

	// Validate product exists
	if order.ProductID != nil {
		_, err := s.productRepo.GetById(*order.ProductID)
		if err != nil {
			return errors.New("product not found")
		}
	}

	// Validate product item exists and is available
	if order.ProductItemID != nil {
		productItem, err := s.productItemRepo.GetByID(*order.ProductItemID)
		if err != nil {
			return errors.New("product item not found")
		}
		if productItem.IsUsed {
			return errors.New("product item is already used")
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
	// Validate status
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

	// Mark product item as used
	if order.ProductItemID != nil {
		if err := s.productItemRepo.MarkAsUsed(*order.ProductItemID); err != nil {
			return err
		}
	}

	// Update order status to complete
	return s.orderRepo.UpdateStatus(orderID, "completed")
}

func (s *orderService) CancelOrder(orderID int, customerID int) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	// Check if user is the customer
	if order.CustomerID == nil || *order.CustomerID != customerID {
		return errors.New("unauthorized to cancel this order")
	}

	if order.Status != "pending" {
		return errors.New("order cannot be cancelled")
	}

	return s.orderRepo.UpdateStatus(orderID, "cancelled")
}

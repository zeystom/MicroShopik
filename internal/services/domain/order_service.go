package domain

import (
	"MicroShopik/internal/domain"
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
}

type orderService struct {
	orderRepo domain.OrderRepository
}

func NewOrderService(oRepo domain.OrderRepository) OrderService {
	return &orderService{
		orderRepo: oRepo,
	}
}

func (s *orderService) Create(order *domain.Order) error {
	if order.Status == "" {
		order.Status = "pending"
	}

	validStatuses := []string{"pending", "confirmed", "completed", "refunded", "cancelled"}
	isValid := false
	for _, validStatus := range validStatuses {
		if order.Status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid order status")
	}

	if order.CustomerID == nil {
		return errors.New("customer ID is required")
	}

	if order.ProductID == nil {
		return errors.New("product ID is required")
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
	validStatuses := []string{"pending", "confirmed", "completed", "refunded", "cancelled"}
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

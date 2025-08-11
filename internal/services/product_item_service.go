package services

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"errors"
)

type ProductItemService interface {
	Create(productItem *domain.ProductItem) error
	GetByID(id int) (*domain.ProductItem, error)
	GetByProductID(productID int) ([]*domain.ProductItem, error)
	GetAvailableByProductID(productID int) ([]*domain.ProductItem, error)
	Update(productItem *domain.ProductItem) error
	Delete(id int) error
	MarkAsUsed(id int) error
	CreateBulk(productID int, dataItems []string) error
}

type productItemService struct {
	productItemRepo repositories.ProductItemRepository
	productRepo     repositories.ProductCRUDRepository
}

func NewProductItemService(piRepo repositories.ProductItemRepository, pRepo repositories.ProductCRUDRepository) ProductItemService {
	return &productItemService{
		productItemRepo: piRepo,
		productRepo:     pRepo,
	}
}

func (s *productItemService) Create(productItem *domain.ProductItem) error {
	// Validate product exists
	_, err := s.productRepo.GetById(productItem.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	return s.productItemRepo.Create(productItem)
}

func (s *productItemService) GetByID(id int) (*domain.ProductItem, error) {
	return s.productItemRepo.GetByID(id)
}

func (s *productItemService) GetByProductID(productID int) ([]*domain.ProductItem, error) {
	return s.productItemRepo.GetByProductID(productID)
}

func (s *productItemService) GetAvailableByProductID(productID int) ([]*domain.ProductItem, error) {
	return s.productItemRepo.GetAvailableByProductID(productID)
}

func (s *productItemService) Update(productItem *domain.ProductItem) error {
	// Validate product exists
	_, err := s.productRepo.GetById(productItem.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	return s.productItemRepo.Update(productItem)
}

func (s *productItemService) Delete(id int) error {
	return s.productItemRepo.Delete(id)
}

func (s *productItemService) MarkAsUsed(id int) error {
	return s.productItemRepo.MarkAsUsed(id)
}

func (s *productItemService) CreateBulk(productID int, dataItems []string) error {
	// Validate product exists
	_, err := s.productRepo.GetById(productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Create multiple product items
	for _, data := range dataItems {
		productItem := &domain.ProductItem{
			ProductID: productID,
			Data:      data,
			IsUsed:    false,
		}
		if err := s.productItemRepo.Create(productItem); err != nil {
			return err
		}
	}

	return nil
}

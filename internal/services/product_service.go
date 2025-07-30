package services

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"errors"
	"time"
)

type ProductService interface {
	Create(p *domain.Product, userID int) (int, error)
	GetById(id int) (*domain.Product, error)
	Update(id int, product *domain.Product, userID int) error
	Delete(id int, userID int) error
	Find(params repositories.ProductQueryParams) ([]*domain.Product, error)
	Count(params repositories.ProductQueryParams) (int, error)
	IsAvailable(id int) (bool, error)
	IncrementSoldCount(id int, delta int) error
}

type productService struct {
	productRepo repositories.ProductCRUDRepository
}

func NewProductService(r repositories.ProductCRUDRepository) ProductService {
	return &productService{productRepo: r}
}

func (s *productService) Create(p *domain.Product, userID int) (int, error) {
	if p == nil {
		return 0, errors.New("product is nil")
	}
	if len(p.Title) <= 0 {
		return 0, errors.New("product title is empty")
	}
	if len(p.Description) <= 0 {
		return 0, errors.New("product description is empty")
	}
	if p.Price <= 0 {
		return 0, errors.New("product price is zero")
	}

	p.SellerID = userID
	p.CreatedAt = time.Now()

	return s.productRepo.Create(p)

}

func (s *productService) GetById(id int) (*domain.Product, error) {
	return s.productRepo.GetById(id)
}

func (s *productService) Update(id int, product *domain.Product, userID int) error {

	if product == nil {
		return errors.New("product is nil")
	}

	existingProduct, err := s.productRepo.GetById(id)
	if err != nil {
		return err
	}
	if existingProduct.SellerID != userID {
		return errors.New("unauthorized: you can only update your own products")
	}

	updateData := repositories.ProductUpdateData{}
	if product.Title != "" {
		updateData.Title = &product.Title
	}
	if product.Description != "" {
		updateData.Description = &product.Description
	}
	if product.Price > 0 {
		updateData.Price = &product.Price
	}
	if product.CategoryID > 0 {
		updateData.CategoryID = &product.CategoryID
	}

	updateData.IsActive = &product.IsActive
	updateData.Disposable = &product.Disposable
	if product.MaxSales > 0 {
		updateData.MaxSales = &product.MaxSales
	}
	return s.productRepo.Update(id, updateData)

}

func (s *productService) Delete(id int, userID int) error {
	existingProduct, err := s.productRepo.GetById(id)
	if err != nil {
		return err
	}
	if existingProduct.SellerID != userID {
		return errors.New("unauthorized: you can only delete your own products")
	}

	return s.productRepo.Delete(id)
}

func (s *productService) IsAvailable(id int) (bool, error) {
	return s.productRepo.IsAvailable(id)
}

func (s *productService) IncrementSoldCount(id int, delta int) error {
	if delta <= 0 {
		return errors.New("delta must be positive")
	}

	available, err := s.IsAvailable(id)
	if err != nil {
		return err
	}
	if !available {
		return errors.New("cannot increment sold count: product is not available")
	}

	product, err := s.GetById(id)
	if err != nil {
		return err
	}

	if product.SoldCount+delta > product.MaxSales {
		return errors.New("cannot increment sold count: would exceed max sales limit")
	}

	return s.productRepo.IncrementSoldCount(id, delta)
}

func (s *productService) Find(params repositories.ProductQueryParams) ([]*domain.Product, error) {
	return s.productRepo.Find(params)
}

func (s *productService) Count(params repositories.ProductQueryParams) (int, error) {
	return s.productRepo.Count(params)
}

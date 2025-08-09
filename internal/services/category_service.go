package services

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"errors"
	"time"
)

type CategoryService interface {
	Create(category *domain.Category) error
	Update(category *domain.Category) error
	Delete(category *domain.Category) error
	GetCategoryById(id int) (*domain.Category, error)
	GetAllCategories() (*[]domain.Category, error)
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(c repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: c}
}

func (c *categoryService) Create(category *domain.Category) error {
	if category == nil {
		return errors.New("category is nil")
	}
	if len(category.Name) <= 0 {
		return errors.New("category name is empty")
	}

	category.CreatedAt = time.Now()
	return c.categoryRepo.Create(category)
}

func (c *categoryService) Update(category *domain.Category) error {
	if category == nil {
		return errors.New("category is nil")
	}
	if len(category.Name) <= 0 {
		return errors.New("category name is empty")
	}
	_, err := c.categoryRepo.Update(category)
	return err
}

func (c *categoryService) Delete(category *domain.Category) error {
	if category == nil {
		return errors.New("category is nil")
	}
	return c.categoryRepo.Delete(category)
}

func (c *categoryService) GetCategoryById(id int) (*domain.Category, error) {
	return c.categoryRepo.GetCategoryById(id)
}
func (c *categoryService) GetAllCategories() (*[]domain.Category, error) {
	return c.categoryRepo.GetAllCategories()
}

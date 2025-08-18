package domain

import (
	"MicroShopik/internal/domain"
	"errors"
	"time"
)

type CategoryService interface {
	Create(category *domain.Category) error
	Update(category *domain.Category) error
	Delete(category *domain.Category) error
	GetCategoryById(id int) (*domain.Category, error)
	GetByID(id int) (*domain.Category, error)
	GetAllCategories() (*[]domain.Category, error)
	ValidateCategoryExists(categoryID int) error
}

type categoryService struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryService(c domain.CategoryRepository) CategoryService {
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

func (c *categoryService) GetByID(id int) (*domain.Category, error) {
	return c.categoryRepo.GetCategoryById(id)
}

func (c *categoryService) ValidateCategoryExists(categoryID int) error {
	_, err := c.categoryRepo.GetCategoryById(categoryID)
	if err != nil {
		return errors.New("category not found")
	}
	return nil
}

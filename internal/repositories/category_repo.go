package repositories

import (
	"MicroShopik/internal/domain"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (c *categoryRepository) Create(category *domain.Category) error {
	return c.db.Create(category).Error
}

func (c *categoryRepository) GetCategoryById(id int) (*domain.Category, error) {
	var category domain.Category
	err := c.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category with id %d not found", id)
		}
		return nil, err
	}
	return &category, err
}

func (c *categoryRepository) GetByID(id int) (*domain.Category, error) {
	return c.GetCategoryById(id)
}

func (c *categoryRepository) GetAllCategories() (*[]domain.Category, error) {
	var categories []domain.Category
	err := c.db.Find(&categories).Error
	return &categories, err
}
func (c *categoryRepository) Update(category *domain.Category) (*domain.Category, error) {
	tx := c.db.Save(category)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, fmt.Errorf("no rows updated for category id %d", category.ID)
	}
	return c.GetCategoryById(category.ID)

}
func (c *categoryRepository) Delete(category *domain.Category) error {

	tx := c.db.Where("id = ?", category.ID).Delete(&domain.Category{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("no category found with id %d to delete", category.ID)
	}
	return nil
}

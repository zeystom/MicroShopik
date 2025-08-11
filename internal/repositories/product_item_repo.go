package repositories

import (
	"MicroShopik/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type productItemRepository struct {
	db *gorm.DB
}

func NewProductItemRepository(db *gorm.DB) ProductItemRepository {
	return &productItemRepository{db: db}
}

func (r *productItemRepository) Create(productItem *domain.ProductItem) error {
	return r.db.Create(productItem).Error
}

func (r *productItemRepository) GetByID(id int) (*domain.ProductItem, error) {
	var productItem domain.ProductItem
	err := r.db.Preload("Product").Where("id = ?", id).First(&productItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product item not found")
		}
		return nil, err
	}
	return &productItem, nil
}

func (r *productItemRepository) GetByProductID(productID int) ([]*domain.ProductItem, error) {
	var productItems []*domain.ProductItem
	err := r.db.Preload("Product").Where("product_id = ?", productID).Find(&productItems).Error
	if err != nil {
		return nil, err
	}
	return productItems, nil
}

func (r *productItemRepository) GetAvailableByProductID(productID int) ([]*domain.ProductItem, error) {
	var productItems []*domain.ProductItem
	err := r.db.Preload("Product").Where("product_id = ? AND is_used = ?", productID, false).Find(&productItems).Error
	if err != nil {
		return nil, err
	}
	return productItems, nil
}

func (r *productItemRepository) Update(productItem *domain.ProductItem) error {
	return r.db.Save(productItem).Error
}

func (r *productItemRepository) Delete(id int) error {
	return r.db.Delete(&domain.ProductItem{}, id).Error
}

func (r *productItemRepository) MarkAsUsed(id int) error {
	return r.db.Model(&domain.ProductItem{}).
		Where("id = ?", id).
		UpdateColumn("is_used", true).Error
}

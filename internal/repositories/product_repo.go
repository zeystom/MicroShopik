package repositories

import (
	"MicroShopik/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(p *domain.Product) (int, error) {
	err := r.db.Create(p).Error
	return p.ID, err
}

func (r *productRepository) GetById(id int) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(id int, data domain.ProductUpdateData) error {
	updates := make(map[string]interface{})

	if data.Title != nil {
		updates["title"] = *data.Title
	}
	if data.Description != nil {
		updates["description"] = *data.Description
	}
	if data.Price != nil {
		updates["price"] = *data.Price
	}
	if data.CategoryID != nil {
		updates["category_id"] = *data.CategoryID
	}
	if data.IsActive != nil {
		updates["is_active"] = *data.IsActive
	}
	if data.Disposable != nil {
		updates["disposable"] = *data.Disposable
	}
	if data.MaxSales != nil {
		updates["max_sales"] = *data.MaxSales
	}

	if len(updates) == 0 {
		return nil // nothing to update
	}

	return r.db.Model(&domain.Product{}).Where("id = ?", id).Updates(updates).Error
}

func (r *productRepository) Delete(id int) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepository) IsAvailable(id int) (bool, error) {
	var product domain.Product
	err := r.db.Select("is_active, max_sales, sold_count").Where("id = ?", id).First(&product).Error
	if err != nil {
		return false, err
	}

	return product.IsActive && (product.MaxSales == 0 || product.SoldCount < product.MaxSales), nil
}

func (r *productRepository) IncrementSoldCount(id int, delta int) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", id).
		UpdateColumn("sold_count", gorm.Expr("sold_count + ?", delta)).Error
}

func (r *productRepository) ReserveProduct(id int) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", id).
		UpdateColumn("sold_count", gorm.Expr("sold_count + ?", 1)).Error
}

func (r *productRepository) CheckAvailabilityAndIncrementSoldCount(id int, delta int) (bool, error) {
	var result struct {
		IsAvailable  bool
		RowsAffected int64
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var product domain.Product
		if err := tx.Select("is_active, max_sales, sold_count").Where("id = ?", id).First(&product).Error; err != nil {
			return err
		}

		isAvailable := product.IsActive && (product.MaxSales == 0 || product.SoldCount < product.MaxSales)
		if !isAvailable {
			result.IsAvailable = false
			return nil
		}

		updateResult := tx.Model(&domain.Product{}).
			Where("id = ? AND is_active = ? AND (max_sales = 0 OR sold_count < max_sales)",
				id, true).
			UpdateColumn("sold_count", gorm.Expr("sold_count + ?", delta))

		if updateResult.Error != nil {
			return updateResult.Error
		}

		result.RowsAffected = updateResult.RowsAffected
		result.IsAvailable = result.RowsAffected > 0

		return nil
	})

	return result.IsAvailable, err
}

func (r *productRepository) Find(params domain.ProductQueryParams) ([]*domain.Product, error) {
	query := r.db.Model(&domain.Product{})

	if params.SellerId != nil {
		query = query.Where("seller_id = ?", *params.SellerId)
	}
	if params.CategoryID != nil {
		query = query.Where("category_id = ?", *params.CategoryID)
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}
	if params.Disposable != nil {
		query = query.Where("disposable = ?", *params.Disposable)
	}
	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}
	if params.SearchQuery != nil && *params.SearchQuery != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?",
			"%"+*params.SearchQuery+"%", "%"+*params.SearchQuery+"%")
	}

	query = query.Order("created_at DESC")

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	var products []*domain.Product
	err := query.Find(&products).Error
	return products, err
}

func (r *productRepository) Count(params domain.ProductQueryParams) (int, error) {
	query := r.db.Model(&domain.Product{})

	if params.SellerId != nil {
		query = query.Where("seller_id = ?", *params.SellerId)
	}
	if params.CategoryID != nil {
		query = query.Where("category_id = ?", *params.CategoryID)
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}
	if params.Disposable != nil {
		query = query.Where("disposable = ?", *params.Disposable)
	}
	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}
	if params.SearchQuery != nil && *params.SearchQuery != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+*params.SearchQuery+"%", "%"+*params.SearchQuery+"%")
	}

	var count int64
	err := query.Count(&count).Error
	return int(count), err
}

func (r *productRepository) GetAll() ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}

func (r *productRepository) GetBySellerID(sellerID int) ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.Preload("Category").Where("seller_id = ?", sellerID).Find(&products).Error
	return products, err
}

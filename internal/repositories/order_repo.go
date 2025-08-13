package repositories

import (
	"MicroShopik/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {

	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id int) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Customer").Preload("Product").Preload("ProductItem").Preload("Messages.Sender").
		Where("id = ?", id).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByCustomerID(customerID int) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.Preload("Product").Preload("ProductItem").
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetByStatus(status string) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.Preload("Customer").Preload("Product").Preload("ProductItem").
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) Delete(id int) error {
	return r.db.Delete(&domain.Order{}, id).Error
}

func (r *orderRepository) UpdateStatus(id int, status string) error {
	return r.db.Model(&domain.Order{}).
		Where("id = ?", id).
		UpdateColumn("status", status).Error
}

func (r *orderRepository) GetByProductID(productID int) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.Preload("Customer").Preload("ProductItem").
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *orderRepository) UpdateStatusTx(tx *gorm.DB, id int, status string) error {
	return tx.Model(&domain.Order{}).
		Where("id = ?", id).
		UpdateColumn("status", status).Error
}

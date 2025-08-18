package domain

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID         int            `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID *int           `json:"customer_id"`
	ProductID  *int           `json:"product_id"`
	Status     string         `json:"status" gorm:"default:'pending';size:20"`
	Customer   *User          `json:"customer" gorm:"foreignKey:CustomerID"`
	Product    *Product       `json:"product" gorm:"foreignKey:ProductID"`
	Messages   []Message      `json:"messages" gorm:"foreignKey:OrderID"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

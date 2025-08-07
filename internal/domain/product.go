package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	SellerID    int            `json:"seller_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null;size:255"`
	Description string         `json:"description" gorm:"type:text"`
	Price       int64          `json:"price" gorm:"not null"`
	CategoryID  int            `json:"category_id" gorm:"not null"`
	Category    Category       `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
	Disposable  bool           `json:"disposable" gorm:"default:false"`
	MaxSales    int            `json:"max_sales" gorm:"not null;default:0"`
	SoldCount   int            `json:"sold_count" gorm:"not null;default:0"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

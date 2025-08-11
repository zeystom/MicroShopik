package domain

import (
	"gorm.io/gorm"
	"time"
)

type ProductItem struct {
	ID        int            `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID int            `json:"product_id" gorm:"not null"`
	Product   Product        `json:"product" gorm:"foreignKey:ProductID"`
	Data      string         `json:"data" gorm:"not null;type:text"`
	IsUsed    bool           `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

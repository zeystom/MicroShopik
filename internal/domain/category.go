package domain

import (
	"time"
)

type Category struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null;uniqueIndex:uni_categories_name;size:50"`
	Products  []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

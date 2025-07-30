package domain

import "time"

type Product struct {
	ID          int       `json:"id"`
	SellerID    int       `json:"seller_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	CategoryID  int       `json:"category_id"`
	IsActive    bool      `json:"is_active"`
	Disposable  bool      `json:"disposable"`
	MaxSales    int       `json:"max_sales"`
	SoldCount   int       `json:"sold_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

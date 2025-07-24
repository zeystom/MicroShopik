package domain

type Product struct {
	ID          int    `json:"id"`
	SellerID    int    `json:"seller_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CategoryID  int    `json:"category_id"`
	IsActive    bool   `json:"is_active"`
	Disposable  bool   `json:"disposable"`
	MaxSales    int    `json:"max_sales"`
	SoldCount   int    `json:"sold_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

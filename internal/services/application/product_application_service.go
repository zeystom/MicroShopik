package application

import (
	"MicroShopik/internal/domain"
	domain2 "MicroShopik/internal/services/domain"
	"errors"
)

type ProductApplicationService struct {
	productService  domain2.ProductService
	categoryService domain2.CategoryService
	userService     domain2.UserService
	orderService    domain2.OrderService
}

func NewProductApplicationService(
	productService domain2.ProductService,
	categoryService domain2.CategoryService,
	userService domain2.UserService,
	orderService domain2.OrderService,
) *ProductApplicationService {
	return &ProductApplicationService{
		productService:  productService,
		categoryService: categoryService,
		userService:     userService,
		orderService:    orderService,
	}
}

func (s *ProductApplicationService) CreateProductWithValidation(product *domain.Product, sellerID int) error {
	if err := s.userService.ValidateUserExists(sellerID); err != nil {
		return errors.New("seller not found")
	}

	if err := s.categoryService.ValidateCategoryExists(product.CategoryID); err != nil {
		return errors.New("category not found")
	}

	product.SellerID = sellerID
	product.IsActive = true
	product.SoldCount = 0

	id, err := s.productService.Create(product, sellerID)
	if err != nil {
		return err
	}
	product.ID = id
	return nil
}

func (s *ProductApplicationService) UpdateProductWithInventoryCheck(product *domain.Product, sellerID int) error {
	existingProduct, err := s.productService.GetById(product.ID)
	if err != nil {
		return errors.New("product not found")
	}

	if existingProduct.SellerID != sellerID {
		return errors.New("unauthorized to update this product")
	}

	if err := s.categoryService.ValidateCategoryExists(product.CategoryID); err != nil {
		return errors.New("category not found")
	}

	product.SellerID = sellerID
	return s.productService.Update(product.ID, product, sellerID)
}

func (s *ProductApplicationService) DeactivateProductWithOrderCheck(productID, sellerID int) error {
	product, err := s.productService.GetById(productID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.SellerID != sellerID {
		return errors.New("unauthorized to deactivate this product")
	}

	orders, err := s.orderService.GetByProductID(productID)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if order.Status == "pending" || order.Status == "confirmed" {
			return errors.New("cannot deactivate product with active orders")
		}
	}

	product.IsActive = false
	return s.productService.Update(productID, product, sellerID)
}

func (s *ProductApplicationService) GetSellerProductsWithStats(sellerID int) ([]*domain.Product, error) {
	if err := s.userService.ValidateUserExists(sellerID); err != nil {
		return nil, errors.New("seller not found")
	}

	products, err := s.productService.Find(domain.ProductQueryParams{SellerId: &sellerID})
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		orders, err := s.orderService.GetByProductID(product.ID)
		if err != nil {
			continue
		}

		totalRevenue := int64(0)
		for _, order := range orders {
			if order.Status == "completed" {
				totalRevenue += product.Price
			}
		}
	}

	return products, nil
}

func (s *ProductApplicationService) GetProductWithFullDetails(productID int) (*domain.Product, error) {
	product, err := s.productService.GetById(productID)
	if err != nil {
		return nil, err
	}

	category, err := s.categoryService.GetByID(product.CategoryID)
	if err != nil {
		return nil, err
	}

	product.Category = *category
	return product, nil
}

func (s *ProductApplicationService) ValidateProductForPurchase(productID, customerID int) error {
	product, err := s.productService.GetById(productID)
	if err != nil {
		return errors.New("product not found")
	}

	if !product.IsActive {
		return errors.New("product is not available")
	}

	if product.Disposable && product.SoldCount >= product.MaxSales {
		return errors.New("product is out of stock")
	}

	if product.SellerID == customerID {
		return errors.New("cannot purchase your own product")
	}

	return nil
}

func (s *ProductApplicationService) FindProducts(params domain.ProductQueryParams) ([]*domain.Product, error) {
	return s.productService.Find(params)
}

func (s *ProductApplicationService) CountProducts(params domain.ProductQueryParams) (int, error) {
	return s.productService.Count(params)
}

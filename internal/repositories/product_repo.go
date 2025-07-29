package repositories

import (
	"MicroShopik/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductCRUDRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(p *domain.Product) (int, error) {
	query := `
		INSERT INTO products (
			seller_id, title, description, price, category_id, is_active,
			disposable, max_sales, sold_count, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,0,NOW(),NOW())
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		p.SellerID, p.Title, p.Description, p.Price, p.CategoryID,
		p.IsActive, p.Disposable, p.MaxSales,
	).Scan(&p.ID)

	return p.ID, err
}

func (r *productRepository) GetById(id int) (*domain.Product, error) {
	query := `SELECT id, seller_id, title, description, price, category_id, is_active, disposable, max_sales, sold_count, created_at, updated_at FROM products WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var p domain.Product
	err := row.Scan(
		&p.ID, &p.SellerID, &p.Title, &p.Description, &p.Price,
		&p.CategoryID, &p.IsActive, &p.Disposable, &p.MaxSales,
		&p.SoldCount, &p.CreatedAt, &p.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("product not found")
	}
	return &p, err
}

func (r *productRepository) Update(id int, data ProductUpdateData) error {
	var sets []string
	var args []interface{}
	i := 1

	if data.Title != nil {
		sets = append(sets, fmt.Sprintf("title = $%d", i))
		args = append(args, *data.Title)
		i++
	}
	if data.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", i))
		args = append(args, *data.Description)
		i++
	}
	if data.Price != nil {
		sets = append(sets, fmt.Sprintf("price = $%d", i))
		args = append(args, *data.Price)
		i++
	}
	if data.CategoryID != nil {
		sets = append(sets, fmt.Sprintf("category_id = $%d", i))
		args = append(args, *data.CategoryID)
		i++
	}
	if data.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", i))
		args = append(args, *data.IsActive)
		i++
	}
	if data.Disposable != nil {
		sets = append(sets, fmt.Sprintf("disposable = $%d", i))
		args = append(args, *data.Disposable)
		i++
	}
	if data.MaxSales != nil {
		sets = append(sets, fmt.Sprintf("max_sales = $%d", i))
		args = append(args, *data.MaxSales)
		i++
	}

	if len(sets) == 0 {
		return nil // nothing to update
	}
	sets = append(sets, "updated_at = NOW()")
	args = append(args, id)
	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d", strings.Join(sets, ", "), i)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *productRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *productRepository) IsAvailable(id int) (bool, error) {
	query := `SELECT is_active, max_sales, sold_count FROM products WHERE id = $1`
	var isActive bool
	var maxSales, soldCount int
	err := r.db.QueryRow(query, id).Scan(&isActive, &maxSales, &soldCount)
	if err != nil {
		return false, err
	}
	if !isActive || soldCount >= maxSales {
		return false, nil
	}
	return true, nil
}

func (r *productRepository) IncrementSoldCount(id int, delta int) error {
	query := `UPDATE products SET sold_count = sold_count + $1 WHERE id = $2`
	_, err := r.db.Exec(query, delta, id)
	return err
}

func (r *productRepository) Find(params ProductQueryParams) ([]*domain.Product, error) {
	var (
		query  = `SELECT id, seller_id, title, description, price, category_id, is_active, disposable, max_sales, sold_count, created_at, updated_at FROM products WHERE 1=1`
		args   []interface{}
		cursor = 1
	)

	if params.SellerId != nil {
		query += fmt.Sprintf(" AND seller_id = $%d", cursor)
		args = append(args, *params.SellerId)
		cursor++
	}
	if params.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", cursor)
		args = append(args, *params.CategoryID)
		cursor++
	}
	if params.MinPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", cursor)
		args = append(args, *params.MinPrice)
		cursor++
	}
	if params.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", cursor)
		args = append(args, *params.MaxPrice)
		cursor++
	}
	if params.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", cursor)
		args = append(args, *params.IsActive)
		cursor++
	}
	if params.Disposable != nil {
		query += fmt.Sprintf(" AND disposable = $%d", cursor)
		args = append(args, *params.Disposable)
		cursor++
	}
	if params.SearchQuery != nil {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", cursor, cursor+1)
		like := "%" + *params.SearchQuery + "%"
		args = append(args, like, like)
		cursor += 2
	}
	query += " ORDER BY created_at DESC"

	if params.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", cursor)
		args = append(args, *params.Limit)
		cursor++
	}
	if params.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", cursor)
		args = append(args, *params.Offset)
		cursor++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(
			&p.ID, &p.SellerID, &p.Title, &p.Description, &p.Price,
			&p.CategoryID, &p.IsActive, &p.Disposable, &p.MaxSales,
			&p.SoldCount, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

func (r *productRepository) Count(params ProductQueryParams) (int, error) {
	var (
		query  = `SELECT COUNT(*) FROM products WHERE 1=1`
		args   []interface{}
		cursor = 1
	)

	if params.SellerId != nil {
		query += fmt.Sprintf(" AND seller_id = $%d", cursor)
		args = append(args, *params.SellerId)
		cursor++
	}
	if params.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", cursor)
		args = append(args, *params.CategoryID)
		cursor++
	}
	if params.MinPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", cursor)
		args = append(args, *params.MinPrice)
		cursor++
	}
	if params.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", cursor)
		args = append(args, *params.MaxPrice)
		cursor++
	}
	if params.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", cursor)
		args = append(args, *params.IsActive)
		cursor++
	}
	if params.Disposable != nil {
		query += fmt.Sprintf(" AND disposable = $%d", cursor)
		args = append(args, *params.Disposable)
		cursor++
	}
	if params.SearchQuery != nil {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", cursor, cursor+1)
		like := "%" + *params.SearchQuery + "%"
		args = append(args, like, like)
		cursor += 2
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

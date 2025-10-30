package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	commonpb "github.com/yourusername/kube-ec/proto/common"
	pb "github.com/yourusername/kube-ec/proto/product"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *pb.Product) error {
	query := `
		INSERT INTO products (id, name, description, price_currency, price_amount, stock_quantity, category, sku, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		product.Id,
		product.Name,
		product.Description,
		product.Price.Currency,
		product.Price.Amount,
		product.StockQuantity,
		product.Category,
		product.Sku,
		product.IsActive,
		now,
		now,
	)
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*pb.Product, error) {
	query := `
		SELECT id, name, description, price_currency, price_amount, stock_quantity, category, sku, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	product := &pb.Product{
		Price:     &commonpb.Money{},
		CreatedAt: &commonpb.Timestamp{},
		UpdatedAt: &commonpb.Timestamp{},
	}

	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.Price.Currency,
		&product.Price.Amount,
		&product.StockQuantity,
		&product.Category,
		&product.Sku,
		&product.IsActive,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, err
	}

	product.CreatedAt.Seconds = createdAt.Unix()
	product.UpdatedAt.Seconds = updatedAt.Unix()

	return product, nil
}

func (r *ProductRepository) List(ctx context.Context, page, pageSize int32, category, searchQuery string) ([]*pb.Product, int32, error) {
	baseQuery := `
		SELECT id, name, description, price_currency, price_amount, stock_quantity, category, sku, is_active, created_at, updated_at
		FROM products
		WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) FROM products WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if category != "" {
		baseQuery += fmt.Sprintf(" AND category = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	if searchQuery != "" {
		baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		countQuery += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%"+searchQuery+"%")
		argIdx++
	}

	// カウント取得
	var totalCount int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (page - 1) * pageSize
	baseQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []*pb.Product{}
	for rows.Next() {
		product := &pb.Product{
			Price:     &commonpb.Money{},
			CreatedAt: &commonpb.Timestamp{},
			UpdatedAt: &commonpb.Timestamp{},
		}
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Price.Currency,
			&product.Price.Amount,
			&product.StockQuantity,
			&product.Category,
			&product.Sku,
			&product.IsActive,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		product.CreatedAt.Seconds = createdAt.Unix()
		product.UpdatedAt.Seconds = updatedAt.Unix()
		products = append(products, product)
	}

	return products, totalCount, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *pb.Product) error {
	query := `
		UPDATE products
		SET name = $2, description = $3, price_currency = $4, price_amount = $5,
		    stock_quantity = $6, category = $7, is_active = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		product.Id,
		product.Name,
		product.Description,
		product.Price.Currency,
		product.Price.Amount,
		product.StockQuantity,
		product.Category,
		product.IsActive,
		time.Now(),
	)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ProductRepository) UpdateStock(ctx context.Context, productID string, quantityChange int32) error {
	query := `
		UPDATE products
		SET stock_quantity = stock_quantity + $2, updated_at = $3
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, productID, quantityChange, time.Now())
	return err
}

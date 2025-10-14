package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/yourusername/kube-ec/proto/order"
	commonpb "github.com/yourusername/kube-ec/proto/common"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *pb.Order) error {
	// OrderItemsをJSONに変換
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	query := `
		INSERT INTO orders (id, user_id, items, total_currency, total_amount, status, payment_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	now := time.Now()
	_, err = r.db.ExecContext(ctx, query,
		order.Id,
		order.UserId,
		itemsJSON,
		order.TotalAmount.Currency,
		order.TotalAmount.Amount,
		order.Status.String(),
		order.PaymentId,
		now,
		now,
	)
	return err
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*pb.Order, error) {
	query := `
		SELECT id, user_id, items, total_currency, total_amount, status, payment_id, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	order := &pb.Order{
		TotalAmount: &commonpb.Money{},
		CreatedAt:   &commonpb.Timestamp{},
		UpdatedAt:   &commonpb.Timestamp{},
	}

	var itemsJSON []byte
	var statusStr string
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.Id,
		&order.UserId,
		&itemsJSON,
		&order.TotalAmount.Currency,
		&order.TotalAmount.Amount,
		&statusStr,
		&order.PaymentId,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, err
	}

	// JSONをOrderItemsに変換
	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	// ステータスの変換
	order.Status = pb.OrderStatus(pb.OrderStatus_value[statusStr])

	order.CreatedAt.Seconds = createdAt.Unix()
	order.UpdatedAt.Seconds = updatedAt.Unix()

	return order, nil
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID string, page, pageSize int32, status pb.OrderStatus) ([]*pb.Order, int32, error) {
	baseQuery := `
		SELECT id, user_id, items, total_currency, total_amount, status, payment_id, created_at, updated_at
		FROM orders
		WHERE user_id = $1
	`
	countQuery := "SELECT COUNT(*) FROM orders WHERE user_id = $1"
	args := []interface{}{userID}
	argIdx := 2

	if status != pb.OrderStatus_ORDER_STATUS_UNSPECIFIED {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status.String())
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

	orders := []*pb.Order{}
	for rows.Next() {
		order := &pb.Order{
			TotalAmount: &commonpb.Money{},
			CreatedAt:   &commonpb.Timestamp{},
			UpdatedAt:   &commonpb.Timestamp{},
		}
		var itemsJSON []byte
		var statusStr string
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&order.Id,
			&order.UserId,
			&itemsJSON,
			&order.TotalAmount.Currency,
			&order.TotalAmount.Amount,
			&statusStr,
			&order.PaymentId,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// JSONをOrderItemsに変換
		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal items: %w", err)
		}

		order.Status = pb.OrderStatus(pb.OrderStatus_value[statusStr])
		order.CreatedAt.Seconds = createdAt.Unix()
		order.UpdatedAt.Seconds = updatedAt.Unix()
		orders = append(orders, order)
	}

	return orders, totalCount, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status pb.OrderStatus) error {
	query := `
		UPDATE orders
		SET status = $2, updated_at = $3
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, status.String(), time.Now())
	return err
}

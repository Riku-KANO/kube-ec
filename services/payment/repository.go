package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pb "github.com/yourusername/kube-ec/proto/payment"
	commonpb "github.com/yourusername/kube-ec/proto/common"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *pb.Payment) error {
	query := `
		INSERT INTO payments (id, order_id, user_id, amount_currency, amount, status, method, transaction_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		payment.Id,
		payment.OrderId,
		payment.UserId,
		payment.Amount.Currency,
		payment.Amount.Amount,
		payment.Status.String(),
		payment.Method.String(),
		payment.TransactionId,
		now,
		now,
	)
	return err
}

func (r *PaymentRepository) GetByID(ctx context.Context, id string) (*pb.Payment, error) {
	query := `
		SELECT id, order_id, user_id, amount_currency, amount, status, method, transaction_id, created_at, updated_at
		FROM payments
		WHERE id = $1
	`
	payment := &pb.Payment{
		Amount:    &commonpb.Money{},
		CreatedAt: &commonpb.Timestamp{},
		UpdatedAt: &commonpb.Timestamp{},
	}

	var statusStr, methodStr string
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.UserId,
		&payment.Amount.Currency,
		&payment.Amount.Amount,
		&statusStr,
		&methodStr,
		&payment.TransactionId,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, err
	}

	payment.Status = pb.PaymentStatus(pb.PaymentStatus_value[statusStr])
	payment.Method = pb.PaymentMethod(pb.PaymentMethod_value[methodStr])
	payment.CreatedAt.Seconds = createdAt.Unix()
	payment.UpdatedAt.Seconds = updatedAt.Unix()

	return payment, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id string, status pb.PaymentStatus, transactionID string) error {
	query := `
		UPDATE payments
		SET status = $2, transaction_id = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, status.String(), transactionID, time.Now())
	return err
}

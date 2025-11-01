package main

import (
	"context"
	"fmt"

	pb "github.com/Riku-KANO/kube-ec/services/payment/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
	repo *PaymentRepository
}

func NewPaymentServer(repo *PaymentRepository) *PaymentServer {
	return &PaymentServer{
		repo: repo,
	}
}

func (s *PaymentServer) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.Payment, error) {
	if req.OrderId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id and user_id are required")
	}
	if req.Amount == nil || req.Amount.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "valid amount is required")
	}

	payment := &pb.Payment{
		Id:      uuid.New().String(),
		OrderId: req.OrderId,
		UserId:  req.UserId,
		Amount:  req.Amount,
		Status:  pb.PaymentStatus_PAYMENT_STATUS_PENDING,
		Method:  req.Method,
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create payment: %v", err))
	}

	return payment, nil
}

func (s *PaymentServer) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.Payment, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	payment, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "payment not found")
	}

	return payment, nil
}

func (s *PaymentServer) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	if req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	payment, err := s.repo.GetByID(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "payment not found")
	}

	if payment.Status != pb.PaymentStatus_PAYMENT_STATUS_PENDING {
		return &pb.ProcessPaymentResponse{
			Success: false,
			Message: "payment is not in pending status",
		}, nil
	}

	// 実際の決済処理をシミュレート
	// 本番環境では、Stripe、PayPalなどの決済APIを呼び出す
	transactionID := fmt.Sprintf("txn_%s", uuid.New().String())

	// ステータスを更新
	if err := s.repo.UpdateStatus(ctx, req.PaymentId, pb.PaymentStatus_PAYMENT_STATUS_COMPLETED, transactionID); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update payment status: %v", err))
	}

	return &pb.ProcessPaymentResponse{
		Success:       true,
		TransactionId: transactionID,
		Message:       "payment processed successfully",
	}, nil
}

func (s *PaymentServer) RefundPayment(ctx context.Context, req *pb.RefundPaymentRequest) (*pb.RefundPaymentResponse, error) {
	if req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	payment, err := s.repo.GetByID(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "payment not found")
	}

	if payment.Status != pb.PaymentStatus_PAYMENT_STATUS_COMPLETED {
		return &pb.RefundPaymentResponse{
			Success: false,
			Message: "payment is not completed, cannot refund",
		}, nil
	}

	// 実際の返金処理をシミュレート
	refundID := fmt.Sprintf("rfnd_%s", uuid.New().String())

	// ステータスを更新
	if err := s.repo.UpdateStatus(ctx, req.PaymentId, pb.PaymentStatus_PAYMENT_STATUS_REFUNDED, payment.TransactionId); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update payment status: %v", err))
	}

	return &pb.RefundPaymentResponse{
		Success:  true,
		RefundId: refundID,
		Message:  "refund processed successfully",
	}, nil
}

func (s *PaymentServer) GetPaymentStatus(ctx context.Context, req *pb.GetPaymentStatusRequest) (*pb.GetPaymentStatusResponse, error) {
	if req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	payment, err := s.repo.GetByID(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "payment not found")
	}

	return &pb.GetPaymentStatusResponse{
		Status:        payment.Status,
		TransactionId: payment.TransactionId,
	}, nil
}

package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	commonpb "github.com/yourusername/kube-ec/proto/common"
	pb "github.com/yourusername/kube-ec/proto/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	repo *OrderRepository
}

func NewOrderServer(repo *OrderRepository) *OrderServer {
	return &OrderServer{
		repo: repo,
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one item is required")
	}

	// 合計金額の計算
	var totalAmount int64
	currency := "JPY"
	for _, item := range req.Items {
		if item.Subtotal == nil {
			item.Subtotal = &commonpb.Money{
				Currency: currency,
				Amount:   item.UnitPrice.Amount * int64(item.Quantity),
			}
		}
		totalAmount += item.Subtotal.Amount
	}

	order := &pb.Order{
		Id:     uuid.New().String(),
		UserId: req.UserId,
		Items:  req.Items,
		TotalAmount: &commonpb.Money{
			Currency: currency,
			Amount:   totalAmount,
		},
		Status:          pb.OrderStatus_ORDER_STATUS_PENDING,
		ShippingAddress: req.ShippingAddress,
		CreatedAt:       &commonpb.Timestamp{},
		UpdatedAt:       &commonpb.Timestamp{},
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create order: %v", err))
	}

	return order, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	order, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	return order, nil
}

func (s *OrderServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	page := req.Pagination.Page
	pageSize := req.Pagination.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	orders, totalCount, err := s.repo.ListByUser(ctx, req.UserId, page, pageSize, req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list orders: %v", err))
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	return &pb.ListOrdersResponse{
		Orders: orders,
		Pagination: &commonpb.PaginationResponse{
			TotalCount:  totalCount,
			TotalPages:  totalPages,
			CurrentPage: page,
		},
	}, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.repo.UpdateStatus(ctx, req.Id, req.Status); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update order status: %v", err))
	}

	order, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get updated order")
	}

	return order, nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.Order, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// 現在の注文を取得
	order, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	// キャンセル可能かチェック
	if order.Status == pb.OrderStatus_ORDER_STATUS_SHIPPED ||
		order.Status == pb.OrderStatus_ORDER_STATUS_DELIVERED ||
		order.Status == pb.OrderStatus_ORDER_STATUS_CANCELLED {
		return nil, status.Error(codes.FailedPrecondition, "order cannot be cancelled")
	}

	// ステータスをキャンセルに更新
	if err := s.repo.UpdateStatus(ctx, req.Id, pb.OrderStatus_ORDER_STATUS_CANCELLED); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to cancel order: %v", err))
	}

	order, err = s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get cancelled order")
	}

	return order, nil
}

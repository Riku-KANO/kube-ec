package main

import (
	"context"
	"fmt"

	commonpb "github.com/Riku-KANO/kube-ec/pkg/proto/common"
	pb "github.com/Riku-KANO/kube-ec/services/product/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer
	repo *ProductRepository
}

func NewProductServer(repo *ProductRepository) *ProductServer {
	return &ProductServer{
		repo: repo,
	}
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Price == nil || req.Price.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "valid price is required")
	}

	product := &pb.Product{
		Id:            uuid.New().String(),
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Category:      req.Category,
		ImageUrls:     req.ImageUrls,
		Sku:           req.Sku,
		IsActive:      true,
		CreatedAt:     &commonpb.Timestamp{},
		UpdatedAt:     &commonpb.Timestamp{},
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create product: %v", err))
	}

	return product, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	product, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return product, nil
}

func (s *ProductServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
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

	products, totalCount, err := s.repo.List(ctx, page, pageSize, req.Category, req.SearchQuery)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list products: %v", err))
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	return &pb.ListProductsResponse{
		Products: products,
		Pagination: &commonpb.PaginationResponse{
			TotalCount:  totalCount,
			TotalPages:  totalPages,
			CurrentPage: page,
		},
	}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	existing, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.Price = req.Price
	existing.StockQuantity = req.StockQuantity
	existing.Category = req.Category
	existing.ImageUrls = req.ImageUrls
	existing.IsActive = req.IsActive

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update product: %v", err))
	}

	return existing, nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.repo.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete product: %v", err))
	}

	return &pb.DeleteProductResponse{Success: true}, nil
}

func (s *ProductServer) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.Product, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	if err := s.repo.UpdateStock(ctx, req.ProductId, req.QuantityChange); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update stock: %v", err))
	}

	product, err := s.repo.GetByID(ctx, req.ProductId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get updated product")
	}

	return product, nil
}

func (s *ProductServer) CheckStock(ctx context.Context, req *pb.CheckStockRequest) (*pb.CheckStockResponse, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	product, err := s.repo.GetByID(ctx, req.ProductId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	available := product.StockQuantity >= req.RequiredQuantity

	return &pb.CheckStockResponse{
		Available:    available,
		CurrentStock: product.StockQuantity,
	}, nil
}

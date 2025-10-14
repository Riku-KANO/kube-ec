package main

import (
	"fmt"
	"log"

	productpb "github.com/yourusername/kube-ec/proto/product"
	userpb "github.com/yourusername/kube-ec/proto/user"
	orderpb "github.com/yourusername/kube-ec/proto/order"
	paymentpb "github.com/yourusername/kube-ec/proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	ProductClient productpb.ProductServiceClient
	UserClient    userpb.UserServiceClient
	OrderClient   orderpb.OrderServiceClient
	PaymentClient paymentpb.PaymentServiceClient
}

func NewGRPCClients(productAddr, userAddr, orderAddr, paymentAddr string) (*GRPCClients, error) {
	// Product service connection
	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	log.Printf("Connected to product service at %s", productAddr)

	// User service connection
	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	log.Printf("Connected to user service at %s", userAddr)

	// Order service connection
	orderConn, err := grpc.NewClient(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}
	log.Printf("Connected to order service at %s", orderAddr)

	// Payment service connection
	paymentConn, err := grpc.NewClient(paymentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}
	log.Printf("Connected to payment service at %s", paymentAddr)

	return &GRPCClients{
		ProductClient: productpb.NewProductServiceClient(productConn),
		UserClient:    userpb.NewUserServiceClient(userConn),
		OrderClient:   orderpb.NewOrderServiceClient(orderConn),
		PaymentClient: paymentpb.NewPaymentServiceClient(paymentConn),
	}, nil
}

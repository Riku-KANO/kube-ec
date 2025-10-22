package grpc

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "github.com/Riku-KANO/kube-ec/proto/user"
)

// ClientConfig holds gRPC client configuration
type ClientConfig struct {
	UserServiceAddr string
}

// Clients holds all gRPC clients
type Clients struct {
	UserClient userpb.UserServiceClient
}

// NewClients creates new gRPC clients
func NewClients(config ClientConfig) (*Clients, error) {
	// Connect to user service
	userConn, err := grpc.Dial(
		config.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	return &Clients{
		UserClient: userpb.NewUserServiceClient(userConn),
	}, nil
}

package grpc

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/Riku-KANO/kube-ec/services/auth/proto"
	userpb "github.com/Riku-KANO/kube-ec/services/user/proto"
)

// ClientConfig holds gRPC client configuration
type ClientConfig struct {
	AuthServiceAddr string
	UserServiceAddr string
}

// Clients holds all gRPC clients and connections
type Clients struct {
	AuthClient authpb.AuthServiceClient
	UserClient userpb.UserServiceClient
	authConn   *grpc.ClientConn
	userConn   *grpc.ClientConn
}

// NewClients creates new gRPC clients
func NewClients(config ClientConfig) (*Clients, error) {
	// Connect to auth service
	authConn, err := grpc.Dial(
		config.AuthServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	// Connect to user service
	userConn, err := grpc.Dial(
		config.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		authConn.Close()
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	return &Clients{
		AuthClient: authpb.NewAuthServiceClient(authConn),
		UserClient: userpb.NewUserServiceClient(userConn),
		authConn:   authConn,
		userConn:   userConn,
	}, nil
}

// Close closes all gRPC connections
func (c *Clients) Close() error {
	var err error
	if c.authConn != nil {
		if closeErr := c.authConn.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if c.userConn != nil {
		if closeErr := c.userConn.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}

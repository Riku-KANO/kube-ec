package grpc

import (
	"context"
	"time"

	userpb "github.com/Riku-KANO/kube-ec/proto/user"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/user"
)

const (
	// Default timeout for gRPC calls
	defaultGRPCTimeout = 10 * time.Second
)

// UserRepository implements user.Repository using gRPC
type UserRepository struct {
	client userpb.UserServiceClient
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(client userpb.UserServiceClient) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// FindByID retrieves a user by ID via gRPC
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	req := &userpb.GetUserRequest{Id: id}

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	resp, err := r.client.GetUser(ctx, req)
	if err != nil {
		return nil, mapGRPCError(err)
	}

	return toDomainUser(resp)
}

// Update updates user information via gRPC
func (r *UserRepository) Update(
	ctx context.Context,
	id string,
	name string,
	phoneNumber *user.PhoneNumber,
) (*user.User, error) {
	req := &userpb.UpdateUserRequest{
		Id:   id,
		Name: name,
	}
	if phoneNumber != nil {
		req.PhoneNumber = phoneNumber.String()
	}

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	resp, err := r.client.UpdateUser(ctx, req)
	if err != nil {
		return nil, mapGRPCError(err)
	}

	return toDomainUser(resp)
}

// Delete removes a user via gRPC
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	req := &userpb.DeleteUserRequest{Id: id}

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	_, err := r.client.DeleteUser(ctx, req)
	if err != nil {
		return mapGRPCError(err)
	}

	return nil
}

// toDomainUser converts protobuf User to domain User
func toDomainUser(pbUser *userpb.User) (*user.User, error) {
	email, err := user.NewEmail(pbUser.Email)
	if err != nil {
		return nil, err
	}

	var phoneNumber *user.PhoneNumber
	if pbUser.PhoneNumber != "" {
		phone, err := user.NewPhoneNumber(pbUser.PhoneNumber)
		if err != nil {
			return nil, err
		}
		phoneNumber = &phone
	}

	return user.NewUser(
		pbUser.Id,
		email,
		pbUser.Name,
		phoneNumber,
		timestampToTime(pbUser.CreatedAt),
		timestampToTime(pbUser.UpdatedAt),
	), nil
}

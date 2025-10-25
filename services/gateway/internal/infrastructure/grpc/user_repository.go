package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Riku-KANO/kube-ec/proto/common"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/errors"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/user"
	userpb "github.com/Riku-KANO/kube-ec/proto/user"
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

// Register creates a new user via gRPC
func (r *UserRepository) Register(
	ctx context.Context,
	email user.Email,
	password string,
	name string,
	phoneNumber *user.PhoneNumber,
) (*user.User, user.AuthTokens, error) {
	req := &userpb.RegisterRequest{
		Email:    email.String(),
		Password: password,
		Name:     name,
	}
	if phoneNumber != nil {
		req.PhoneNumber = phoneNumber.String()
	}

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	resp, err := r.client.Register(ctx, req)
	if err != nil {
		return nil, user.AuthTokens{}, mapGRPCError(err)
	}

	domainUser, err := toDomainUser(resp.User)
	if err != nil {
		return nil, user.AuthTokens{}, err
	}

	tokens := user.NewAuthTokens(resp.AccessToken, resp.RefreshToken)

	return domainUser, tokens, nil
}

// Login authenticates a user via gRPC
func (r *UserRepository) Login(
	ctx context.Context,
	email user.Email,
	password string,
) (*user.User, user.AuthTokens, error) {
	req := &userpb.LoginRequest{
		Email:    email.String(),
		Password: password,
	}

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	resp, err := r.client.Login(ctx, req)
	if err != nil {
		return nil, user.AuthTokens{}, mapGRPCError(err)
	}

	domainUser, err := toDomainUser(resp.User)
	if err != nil {
		return nil, user.AuthTokens{}, err
	}

	tokens := user.NewAuthTokens(resp.AccessToken, resp.RefreshToken)

	return domainUser, tokens, nil
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

// timestampToTime converts common.Timestamp to time.Time
func timestampToTime(ts *common.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return time.Unix(ts.Seconds, int64(ts.Nanos))
}

// mapGRPCError maps gRPC errors to domain errors with better context
func mapGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return errors.ErrInternalError
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return errors.ErrInvalidInput
	case codes.NotFound:
		return errors.ErrNotFound
	case codes.AlreadyExists:
		return errors.ErrInvalidInput
	case codes.Unauthenticated:
		return errors.ErrUnauthorized
	case codes.PermissionDenied:
		return errors.ErrUnauthorized
	case codes.DeadlineExceeded:
		return errors.ErrInternalError
	case codes.Unavailable:
		return errors.ErrInternalError
	default:
		return errors.ErrInternalError
	}
}

// withTimeout creates a context with timeout for gRPC calls
func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		// Context already has deadline, don't override
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, defaultGRPCTimeout)
}

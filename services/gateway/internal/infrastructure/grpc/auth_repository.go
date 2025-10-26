package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "github.com/Riku-KANO/kube-ec/proto/auth"
	"github.com/Riku-KANO/kube-ec/proto/common"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/errors"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/user"
)

// AuthRepository implements authentication operations using gRPC
type AuthRepository struct {
	client authpb.AuthServiceClient
}

// NewAuthRepository creates a new AuthRepository
func NewAuthRepository(client authpb.AuthServiceClient) *AuthRepository {
	return &AuthRepository{
		client: client,
	}
}

// Register creates a new user via auth service
func (r *AuthRepository) Register(
	ctx context.Context,
	email user.Email,
	password string,
	name string,
	phoneNumber *user.PhoneNumber,
) (*user.User, user.AuthTokens, error) {
	req := &authpb.RegisterRequest{
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

	domainUser, err := authResponseToDomainUser(resp)
	if err != nil {
		return nil, user.AuthTokens{}, err
	}

	tokens := user.NewAuthTokens(resp.AccessToken, resp.RefreshToken)

	return domainUser, tokens, nil
}

// Login authenticates a user via auth service
func (r *AuthRepository) Login(
	ctx context.Context,
	email user.Email,
	password string,
) (*user.User, user.AuthTokens, error) {
	req := &authpb.LoginRequest{
		Email:    email.String(),
		Password: password,
	}

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	resp, err := r.client.Login(ctx, req)
	if err != nil {
		return nil, user.AuthTokens{}, mapGRPCError(err)
	}

	domainUser, err := loginResponseToDomainUser(resp)
	if err != nil {
		return nil, user.AuthTokens{}, err
	}

	tokens := user.NewAuthTokens(resp.AccessToken, resp.RefreshToken)

	return domainUser, tokens, nil
}

// VerifyToken verifies a JWT token via auth service
func (r *AuthRepository) VerifyToken(ctx context.Context, token string) (bool, string, error) {
	req := &authpb.VerifyTokenRequest{
		Token: token,
	}

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	resp, err := r.client.VerifyToken(ctx, req)
	if err != nil {
		return false, "", mapGRPCError(err)
	}

	return resp.Valid, resp.UserId, nil
}

// RefreshToken generates new tokens via auth service
func (r *AuthRepository) RefreshToken(ctx context.Context, refreshToken string) (user.AuthTokens, error) {
	req := &authpb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	resp, err := r.client.RefreshToken(ctx, req)
	if err != nil {
		return user.AuthTokens{}, mapGRPCError(err)
	}

	tokens := user.NewAuthTokens(resp.AccessToken, resp.RefreshToken)
	return tokens, nil
}

// authResponseToDomainUser converts auth RegisterResponse to domain User
func authResponseToDomainUser(resp *authpb.RegisterResponse) (*user.User, error) {
	email, err := user.NewEmail(resp.Email)
	if err != nil {
		return nil, err
	}

	var phoneNumber *user.PhoneNumber
	if resp.PhoneNumber != "" {
		phone, err := user.NewPhoneNumber(resp.PhoneNumber)
		if err != nil {
			return nil, err
		}
		phoneNumber = &phone
	}

	return user.NewUser(
		resp.UserId,
		email,
		resp.Name,
		phoneNumber,
		timestampToTime(resp.CreatedAt),
		time.Now(), // UpdatedAt is same as CreatedAt for new users
	), nil
}

// loginResponseToDomainUser converts auth LoginResponse to domain User
func loginResponseToDomainUser(resp *authpb.LoginResponse) (*user.User, error) {
	email, err := user.NewEmail(resp.Email)
	if err != nil {
		return nil, err
	}

	var phoneNumber *user.PhoneNumber
	if resp.PhoneNumber != "" {
		phone, err := user.NewPhoneNumber(resp.PhoneNumber)
		if err != nil {
			return nil, err
		}
		phoneNumber = &phone
	}

	return user.NewUser(
		resp.UserId,
		email,
		resp.Name,
		phoneNumber,
		time.Now(), // We don't have timestamps in login response, use current time
		time.Now(),
	), nil
}

// timestampToTime converts common.Timestamp to time.Time
func timestampToTime(ts *common.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return time.Unix(ts.Seconds, int64(ts.Nanos))
}

// mapGRPCError maps gRPC errors to domain errors
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

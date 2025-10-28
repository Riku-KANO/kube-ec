package grpc

import (
	"context"

	authpb "github.com/Riku-KANO/kube-ec/proto/auth"
	commonpb "github.com/Riku-KANO/kube-ec/proto/common"
	pkgerrors "github.com/Riku-KANO/kube-ec/pkg/errors"
	appauth "github.com/Riku-KANO/kube-ec/services/auth/internal/application/auth"
)

// AuthHandler implements the gRPC AuthService
type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	authService *appauth.Service
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *appauth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register creates a new user account
func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	// Validate request
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, pkgerrors.ErrInvalidArgument.GRPCStatus().Err()
	}

	// Call application service
	output, err := h.authService.Register(ctx, appauth.RegisterInput{
		Email:       req.Email,
		Password:    req.Password,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	})

	if err != nil {
		if pkgErr, ok := err.(*pkgerrors.Error); ok {
			return nil, pkgErr.GRPCStatus().Err()
		}
		return nil, pkgerrors.ErrInternal.GRPCStatus().Err()
	}

	return &authpb.RegisterResponse{
		UserId:       output.UserID,
		Email:        output.Email,
		Name:         output.Name,
		PhoneNumber:  output.PhoneNumber,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		CreatedAt: &commonpb.Timestamp{
			Seconds: output.CreatedAt.Unix(),
		},
	}, nil
}

// Login authenticates a user
func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	// Validate request
	if req.Email == "" || req.Password == "" {
		return nil, pkgerrors.ErrInvalidArgument.GRPCStatus().Err()
	}

	// Call application service
	output, err := h.authService.Login(ctx, appauth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if pkgErr, ok := err.(*pkgerrors.Error); ok {
			return nil, pkgErr.GRPCStatus().Err()
		}
		return nil, pkgerrors.ErrInternal.GRPCStatus().Err()
	}

	return &authpb.LoginResponse{
		UserId:       output.UserID,
		Email:        output.Email,
		Name:         output.Name,
		PhoneNumber:  output.PhoneNumber,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

// VerifyToken validates a JWT token
func (h *AuthHandler) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	if req.Token == "" {
		return nil, pkgerrors.ErrInvalidArgument.GRPCStatus().Err()
	}

	output, err := h.authService.VerifyToken(ctx, req.Token)
	if err != nil {
		return &authpb.VerifyTokenResponse{Valid: false}, nil
	}

	return &authpb.VerifyTokenResponse{
		Valid:  output.Valid,
		UserId: output.UserID,
		Email:  output.Email,
		ExpiresAt: &commonpb.Timestamp{
			Seconds: output.ExpiresAt.Unix(),
		},
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (h *AuthHandler) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, pkgerrors.ErrInvalidArgument.GRPCStatus().Err()
	}

	output, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if pkgErr, ok := err.(*pkgerrors.Error); ok {
			return nil, pkgErr.GRPCStatus().Err()
		}
		return nil, pkgerrors.ErrInternal.GRPCStatus().Err()
	}

	return &authpb.RefreshTokenResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		ExpiresAt: &commonpb.Timestamp{
			Seconds: output.AccessTokenExpiresAt.Unix(),
		},
	}, nil
}

// Logout invalidates a user's tokens
func (h *AuthHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	// TODO: Implement token blacklist/invalidation
	// For now, just return success (client-side token deletion)
	return &authpb.LogoutResponse{Success: true}, nil
}

// ChangePassword changes a user's password
func (h *AuthHandler) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.ChangePasswordResponse, error) {
	if req.UserId == "" || req.OldPassword == "" || req.NewPassword == "" {
		return nil, pkgerrors.ErrInvalidArgument.GRPCStatus().Err()
	}

	err := h.authService.ChangePassword(ctx, appauth.ChangePasswordInput{
		UserID:      req.UserId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		if pkgErr, ok := err.(*pkgerrors.Error); ok {
			return nil, pkgErr.GRPCStatus().Err()
		}
		return nil, pkgerrors.ErrInternal.GRPCStatus().Err()
	}

	return &authpb.ChangePasswordResponse{Success: true}, nil
}

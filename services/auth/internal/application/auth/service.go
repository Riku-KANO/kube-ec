package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/Riku-KANO/kube-ec/pkg/auth"
	pkgerrors "github.com/Riku-KANO/kube-ec/pkg/errors"
	domainauth "github.com/Riku-KANO/kube-ec/services/auth/internal/domain/auth"
)

// Service handles authentication business logic
type Service struct {
	repo            domainauth.Repository
	jwtSecret       string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

// NewService creates a new authentication service
func NewService(
	repo domainauth.Repository,
	jwtSecret string,
	accessDuration time.Duration,
	refreshDuration time.Duration,
) *Service {
	return &Service{
		repo:            repo,
		jwtSecret:       jwtSecret,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, input RegisterInput) (*AuthOutput, error) {
	// Validate email
	email, err := domainauth.NewEmail(input.Email)
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInvalidArgument, err.Error())
	}

	// Validate and hash password
	password, err := domainauth.NewPassword(input.Password)
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInvalidArgument, err.Error())
	}

	// Validate name
	if input.Name == "" {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInvalidArgument, "name is required")
	}

	// Create user entity
	now := time.Now()
	user := domainauth.NewUser(
		uuid.New().String(),
		email,
		password,
		input.Name,
		input.PhoneNumber,
		now,
		now,
	)

	// Save to repository
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessManager := auth.NewJWTManager(s.jwtSecret, s.accessDuration)
	accessToken, err := accessManager.Generate(user.ID(), user.Email().String())
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInternal, "failed to generate access token")
	}

	refreshManager := auth.NewJWTManager(s.jwtSecret, s.refreshDuration)
	refreshToken, err := refreshManager.Generate(user.ID(), user.Email().String())
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInternal, "failed to generate refresh token")
	}

	return &AuthOutput{
		UserID:              user.ID(),
		Email:               user.Email().String(),
		Name:                user.Name(),
		PhoneNumber:         user.PhoneNumber(),
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		AccessTokenExpiresAt: time.Now().Add(s.accessDuration),
		CreatedAt:           user.CreatedAt(),
	}, nil
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, input LoginInput) (*AuthOutput, error) {
	// Validate email
	email, err := domainauth.NewEmail(input.Email)
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInvalidArgument, err.Error())
	}

	// Find user by email
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, pkgerrors.ErrUnauthenticated
	}

	// Verify password
	if err := user.VerifyPassword(input.Password); err != nil {
		return nil, pkgerrors.ErrUnauthenticated
	}

	// Generate tokens
	accessManager := auth.NewJWTManager(s.jwtSecret, s.accessDuration)
	accessToken, err := accessManager.Generate(user.ID(), user.Email().String())
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInternal, "failed to generate access token")
	}

	refreshManager := auth.NewJWTManager(s.jwtSecret, s.refreshDuration)
	refreshToken, err := refreshManager.Generate(user.ID(), user.Email().String())
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInternal, "failed to generate refresh token")
	}

	return &AuthOutput{
		UserID:              user.ID(),
		Email:               user.Email().String(),
		Name:                user.Name(),
		PhoneNumber:         user.PhoneNumber(),
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		AccessTokenExpiresAt: time.Now().Add(s.accessDuration),
		CreatedAt:           user.CreatedAt(),
	}, nil
}

// VerifyToken validates a JWT token
func (s *Service) VerifyToken(ctx context.Context, token string) (*TokenVerificationOutput, error) {
	jwtManager := auth.NewJWTManager(s.jwtSecret, s.accessDuration)
	claims, err := jwtManager.Verify(token)
	if err != nil {
		return &TokenVerificationOutput{Valid: false}, nil
	}

	return &TokenVerificationOutput{
		Valid:     true,
		UserID:    claims.UserID,
		Email:     claims.Email,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*AuthOutput, error) {
	refreshManager := auth.NewJWTManager(s.jwtSecret, s.refreshDuration)
	claims, err := refreshManager.Verify(refreshToken)
	if err != nil {
		return nil, pkgerrors.ErrUnauthenticated
	}

	// Get user to ensure they still exist
	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, pkgerrors.ErrUnauthenticated
	}

	// Generate new tokens
	accessManager := auth.NewJWTManager(s.jwtSecret, s.accessDuration)
	accessToken, err := accessManager.Generate(user.ID(), user.Email().String())
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInternal, "failed to generate access token")
	}

	newRefreshToken, err := refreshManager.Generate(user.ID(), user.Email().String())
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInternal, "failed to generate refresh token")
	}

	return &AuthOutput{
		UserID:              user.ID(),
		Email:               user.Email().String(),
		Name:                user.Name(),
		PhoneNumber:         user.PhoneNumber(),
		AccessToken:         accessToken,
		RefreshToken:        newRefreshToken,
		AccessTokenExpiresAt: time.Now().Add(s.accessDuration),
		CreatedAt:           user.CreatedAt(),
	}, nil
}

// ChangePassword changes a user's password
func (s *Service) ChangePassword(ctx context.Context, input ChangePasswordInput) error {
	// Get user
	user, err := s.repo.FindByID(ctx, input.UserID)
	if err != nil {
		return pkgerrors.ErrNotFound
	}

	// Change password (validates old password internally)
	if err := user.ChangePassword(input.OldPassword, input.NewPassword); err != nil {
		return pkgerrors.Wrap(pkgerrors.ErrUnauthenticated, "invalid old password")
	}

	// Update in repository
	if err := s.repo.UpdatePassword(ctx, user.ID(), user.Password()); err != nil {
		return err
	}

	return nil
}

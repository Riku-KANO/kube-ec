package user

import (
	"context"

	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/errors"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/user"
)

// Service ユーザーアプリケーションサービス
type Service struct {
	userRepo user.Repository
}

// NewService creates a new user application service
func NewService(userRepo user.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

// Register registers a new user
func (s *Service) Register(ctx context.Context, input RegisterInput) (AuthOutput, error) {
	// Validate and create value objects
	email, err := user.NewEmail(input.Email)
	if err != nil {
		return AuthOutput{}, errors.ErrInvalidInput
	}

	if input.Password == "" || input.Name == "" {
		return AuthOutput{}, errors.ErrInvalidInput
	}

	var phoneNumber *user.PhoneNumber
	if input.PhoneNumber != nil && *input.PhoneNumber != "" {
		phone, err := user.NewPhoneNumber(*input.PhoneNumber)
		if err != nil {
			return AuthOutput{}, errors.ErrInvalidInput
		}
		phoneNumber = &phone
	}

	// Call repository
	registeredUser, tokens, err := s.userRepo.Register(ctx, email, input.Password, input.Name, phoneNumber)
	if err != nil {
		return AuthOutput{}, err
	}

	return ToAuthOutput(registeredUser, tokens), nil
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, input LoginInput) (AuthOutput, error) {
	email, err := user.NewEmail(input.Email)
	if err != nil {
		return AuthOutput{}, errors.ErrInvalidInput
	}

	if input.Password == "" {
		return AuthOutput{}, errors.ErrInvalidInput
	}

	authenticatedUser, tokens, err := s.userRepo.Login(ctx, email, input.Password)
	if err != nil {
		return AuthOutput{}, err
	}

	return ToAuthOutput(authenticatedUser, tokens), nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, userID string) (UserOutput, error) {
	if userID == "" {
		return UserOutput{}, errors.ErrInvalidInput
	}

	foundUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return UserOutput{}, err
	}

	return ToUserOutput(foundUser), nil
}

// UpdateUser updates user information
func (s *Service) UpdateUser(ctx context.Context, userID string, input UpdateUserInput) (UserOutput, error) {
	if userID == "" {
		return UserOutput{}, errors.ErrInvalidInput
	}

	var name string
	if input.Name != nil {
		name = *input.Name
	}

	var phoneNumber *user.PhoneNumber
	if input.PhoneNumber != nil && *input.PhoneNumber != "" {
		phone, err := user.NewPhoneNumber(*input.PhoneNumber)
		if err != nil {
			return UserOutput{}, errors.ErrInvalidInput
		}
		phoneNumber = &phone
	}

	updatedUser, err := s.userRepo.Update(ctx, userID, name, phoneNumber)
	if err != nil {
		return UserOutput{}, err
	}

	return ToUserOutput(updatedUser), nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.ErrInvalidInput
	}

	return s.userRepo.Delete(ctx, userID)
}

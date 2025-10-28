package user

import "context"

// AuthRepository defines the interface for authentication operations
type AuthRepository interface {
	// Register creates a new user and returns auth tokens
	Register(ctx context.Context, email Email, password string, name string, phoneNumber *PhoneNumber) (*User, AuthTokens, error)

	// Login authenticates a user and returns auth tokens
	Login(ctx context.Context, email Email, password string) (*User, AuthTokens, error)
}

// UserRepository defines the interface for user management operations
type UserRepository interface {
	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id string) (*User, error)

	// Update updates user information
	Update(ctx context.Context, id string, name string, phoneNumber *PhoneNumber) (*User, error)

	// Delete removes a user
	Delete(ctx context.Context, id string) error
}

package auth

import (
	"context"
)

// Repository defines the interface for authentication data operations
type Repository interface {
	// CreateUser creates a new user in the database
	CreateUser(ctx context.Context, user *User) error

	// FindByEmail retrieves a user by email address
	FindByEmail(ctx context.Context, email Email) (*User, error)

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id string) (*User, error)

	// UpdatePassword updates a user's password
	UpdatePassword(ctx context.Context, userID string, password Password) error
}

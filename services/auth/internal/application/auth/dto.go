package auth

import "time"

// RegisterInput represents registration request data
type RegisterInput struct {
	Email       string
	Password    string
	Name        string
	PhoneNumber string
}

// LoginInput represents login request data
type LoginInput struct {
	Email    string
	Password string
}

// ChangePasswordInput represents password change request data
type ChangePasswordInput struct {
	UserID      string
	OldPassword string
	NewPassword string
}

// AuthOutput represents authentication response data
type AuthOutput struct {
	UserID       string
	Email        string
	Name         string
	PhoneNumber  string
	AccessToken  string
	RefreshToken string
	CreatedAt    time.Time
}

// TokenVerificationOutput represents token verification result
type TokenVerificationOutput struct {
	Valid     bool
	UserID    string
	Email     string
	ExpiresAt time.Time
}

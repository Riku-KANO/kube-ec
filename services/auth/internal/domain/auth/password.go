package auth

import (
	"fmt"

	pkgauth "github.com/Riku-KANO/kube-ec/pkg/auth"
)

// Password represents a hashed password value object
type Password struct {
	hash string
}

const (
	minPasswordLength = 8
	maxPasswordLength = 72 // bcrypt limitation
)

// NewPassword creates a new Password value object from a plain text password
func NewPassword(plainPassword string) (Password, error) {
	// Validate password
	if len(plainPassword) < minPasswordLength {
		return Password{}, fmt.Errorf("password must be at least %d characters", minPasswordLength)
	}
	if len(plainPassword) > maxPasswordLength {
		return Password{}, fmt.Errorf("password must be at most %d characters", maxPasswordLength)
	}

	// Hash password using pkg/auth
	hash, err := pkgauth.HashPassword(plainPassword)
	if err != nil {
		return Password{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return Password{hash: hash}, nil
}

// NewPasswordFromHash creates a Password from an existing hash
func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

// Hash returns the password hash
func (p Password) Hash() string {
	return p.hash
}

// Verify checks if the provided plain password matches the hash
func (p Password) Verify(plainPassword string) error {
	err := pkgauth.VerifyPassword(p.hash, plainPassword)
	if err != nil {
		return fmt.Errorf("invalid password")
	}
	return nil
}

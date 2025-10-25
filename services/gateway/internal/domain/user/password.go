package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Password represents a hashed password value object
type Password struct {
	hash string
}

// NewPassword creates a new Password from a plain text password
func NewPassword(plainPassword string) (Password, error) {
	if len(plainPassword) < 8 {
		return Password{}, fmt.Errorf("password must be at least 8 characters")
	}
	if len(plainPassword) > 72 {
		return Password{}, fmt.Errorf("password must not exceed 72 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return Password{hash: string(hash)}, nil
}

// NewPasswordFromHash creates a Password from an already hashed password
func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

// Hash returns the password hash
func (p Password) Hash() string {
	return p.hash
}

// Verify checks if the plain password matches the hash
func (p Password) Verify(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plainPassword))
	return err == nil
}

package auth

import (
	"fmt"
	"regexp"
	"strings"
)

// Email represents an email address value object
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email value object with validation
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return Email{}, fmt.Errorf("email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return Email{}, fmt.Errorf("invalid email format")
	}

	return Email{value: email}, nil
}

// String returns the email as a string
func (e Email) String() string {
	return e.value
}

// Equals checks if two emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// AuthTokens represents authentication tokens
type AuthTokens struct {
	accessToken  string
	refreshToken string
}

// NewAuthTokens creates a new AuthTokens value object
func NewAuthTokens(accessToken, refreshToken string) AuthTokens {
	return AuthTokens{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

// AccessToken returns the access token
func (t AuthTokens) AccessToken() string {
	return t.accessToken
}

// RefreshToken returns the refresh token
func (t AuthTokens) RefreshToken() string {
	return t.refreshToken
}

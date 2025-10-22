package user

import (
	"fmt"
	"regexp"
)

// Email 値オブジェクト
type Email struct {
	value string
}

// NewEmail creates a new Email value object
func NewEmail(email string) (Email, error) {
	if !isValidEmail(email) {
		return Email{}, fmt.Errorf("invalid email format: %s", email)
	}
	return Email{value: email}, nil
}

func (e Email) String() string {
	return e.value
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// PhoneNumber 値オブジェクト
type PhoneNumber struct {
	value string
}

// NewPhoneNumber creates a new PhoneNumber value object
func NewPhoneNumber(phone string) (PhoneNumber, error) {
	if phone == "" {
		return PhoneNumber{}, fmt.Errorf("phone number cannot be empty")
	}
	// 簡易的なバリデーション（必要に応じて強化）
	return PhoneNumber{value: phone}, nil
}

func (p PhoneNumber) String() string {
	return p.value
}

// AuthTokens 認証トークンの値オブジェクト
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

// NewAuthTokens creates a new AuthTokens value object
func NewAuthTokens(accessToken, refreshToken string) AuthTokens {
	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

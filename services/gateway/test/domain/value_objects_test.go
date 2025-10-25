package domain

import (
	"testing"

	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/user"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "invalid email - no @",
			email:   "testexample.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no domain",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "invalid email - no TLD",
			email:   "test@example",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.email {
				t.Errorf("NewEmail().String() = %v, want %v", got.String(), tt.email)
			}
		})
	}
}

func TestNewPhoneNumber(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{
			name:    "valid Japan phone number",
			phone:   "+819012345678",
			wantErr: false,
		},
		{
			name:    "valid US phone number",
			phone:   "+14155552671",
			wantErr: false,
		},
		{
			name:    "valid UK phone number",
			phone:   "+442071838750",
			wantErr: false,
		},
		{
			name:    "invalid - no plus sign",
			phone:   "819012345678",
			wantErr: true,
		},
		{
			name:    "invalid - starts with 0",
			phone:   "+019012345678",
			wantErr: true,
		},
		{
			name:    "minimum valid length",
			phone:   "+12",
			wantErr: false,
		},
		{
			name:    "invalid - too long",
			phone:   "+81901234567890123",
			wantErr: true,
		},
		{
			name:    "invalid - contains letters",
			phone:   "+81abc12345678",
			wantErr: true,
		},
		{
			name:    "empty phone number",
			phone:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewPhoneNumber(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.phone {
				t.Errorf("NewPhoneNumber().String() = %v, want %v", got.String(), tt.phone)
			}
		})
	}
}

func TestNewAuthTokens(t *testing.T) {
	accessToken := "access-token-123"
	refreshToken := "refresh-token-456"

	tokens := user.NewAuthTokens(accessToken, refreshToken)

	if tokens.AccessToken != accessToken {
		t.Errorf("AuthTokens.AccessToken = %v, want %v", tokens.AccessToken, accessToken)
	}
	if tokens.RefreshToken != refreshToken {
		t.Errorf("AuthTokens.RefreshToken = %v, want %v", tokens.RefreshToken, refreshToken)
	}
}

package domain

import (
	"testing"

	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/user"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name          string
		plainPassword string
		wantErr       bool
	}{
		{
			name:          "valid password",
			plainPassword: "password123",
			wantErr:       false,
		},
		{
			name:          "password too short",
			plainPassword: "pass",
			wantErr:       true,
		},
		{
			name:          "password too long",
			plainPassword: string(make([]byte, 73)),
			wantErr:       true,
		},
		{
			name:          "minimum length password",
			plainPassword: "passwor1",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewPassword(tt.plainPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Hash() == "" {
				t.Error("NewPassword() returned empty hash")
			}
		})
	}
}

func TestPassword_Verify(t *testing.T) {
	plainPassword := "mysecretpassword"
	password, err := user.NewPassword(plainPassword)
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	tests := []struct {
		name          string
		plainPassword string
		want          bool
	}{
		{
			name:          "correct password",
			plainPassword: "mysecretpassword",
			want:          true,
		},
		{
			name:          "incorrect password",
			plainPassword: "wrongpassword",
			want:          false,
		},
		{
			name:          "empty password",
			plainPassword: "",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := password.Verify(tt.plainPassword); got != tt.want {
				t.Errorf("Password.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

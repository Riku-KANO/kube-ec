package user

import (
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/user"
)

// ToUserOutput converts domain User to UserOutput DTO
func ToUserOutput(u *user.User) UserOutput {
	output := UserOutput{
		ID:        u.ID(),
		Email:     u.Email().String(),
		Name:      u.Name(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}

	if u.PhoneNumber() != nil {
		phone := u.PhoneNumber().String()
		output.PhoneNumber = &phone
	}

	return output
}

// ToAuthOutput converts domain User and AuthTokens to AuthOutput DTO
func ToAuthOutput(u *user.User, tokens user.AuthTokens) AuthOutput {
	return AuthOutput{
		User:         ToUserOutput(u),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

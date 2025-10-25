package handler

import (
	appuser "github.com/Riku-KANO/kube-ec/services/gateway/internal/application/user"
	api "github.com/Riku-KANO/kube-ec/services/gateway/internal/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// toRegisterInput converts OpenAPI RegisterRequest to application DTO
func toRegisterInput(req api.RegisterRequest) appuser.RegisterInput {
	input := appuser.RegisterInput{
		Email:    string(req.Email),
		Password: req.Password,
		Name:     req.Name,
	}
	if req.PhoneNumber != nil {
		input.PhoneNumber = req.PhoneNumber
	}
	return input
}

// toLoginInput converts OpenAPI LoginRequest to application DTO
func toLoginInput(req api.LoginRequest) appuser.LoginInput {
	return appuser.LoginInput{
		Email:    string(req.Email),
		Password: req.Password,
	}
}

// toUpdateUserInput converts OpenAPI UpdateUserRequest to application DTO
func toUpdateUserInput(req api.UpdateUserRequest) appuser.UpdateUserInput {
	return appuser.UpdateUserInput{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}
}

// toUserResponse converts application UserOutput to OpenAPI User
func toUserResponse(output appuser.UserOutput) api.User {
	user := api.User{
		Id:        output.ID,
		Email:     openapi_types.Email(output.Email),
		Name:      output.Name,
		CreatedAt: output.CreatedAt,
		UpdatedAt: output.UpdatedAt,
	}
	if output.PhoneNumber != nil {
		user.PhoneNumber = output.PhoneNumber
	}
	return user
}

// toRegisterResponse converts application AuthOutput to OpenAPI RegisterResponse
func toRegisterResponse(output appuser.AuthOutput) api.RegisterResponse {
	return api.RegisterResponse{
		User:         toUserResponse(output.User),
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}
}

// toLoginResponse converts application AuthOutput to OpenAPI LoginResponse
func toLoginResponse(output appuser.AuthOutput) api.LoginResponse {
	return api.LoginResponse{
		User:         toUserResponse(output.User),
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}
}

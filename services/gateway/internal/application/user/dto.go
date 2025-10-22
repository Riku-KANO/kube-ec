package user

import "time"

// RegisterInput ユーザー登録の入力DTO
type RegisterInput struct {
	Email       string
	Password    string
	Name        string
	PhoneNumber *string
}

// LoginInput ログインの入力DTO
type LoginInput struct {
	Email    string
	Password string
}

// UpdateUserInput ユーザー更新の入力DTO
type UpdateUserInput struct {
	Name        *string
	PhoneNumber *string
}

// UserOutput ユーザー情報の出力DTO
type UserOutput struct {
	ID          string
	Email       string
	Name        string
	PhoneNumber *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AuthOutput 認証結果の出力DTO
type AuthOutput struct {
	User         UserOutput
	AccessToken  string
	RefreshToken string
}

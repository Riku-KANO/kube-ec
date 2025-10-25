package user

import "time"

// User ドメインエンティティ
type User struct {
	id          string
	email       Email
	name        string
	phoneNumber *PhoneNumber
	createdAt   time.Time
	updatedAt   time.Time
}

// NewUser creates a new User entity
func NewUser(
	id string,
	email Email,
	name string,
	phoneNumber *PhoneNumber,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:          id,
		email:       email,
		name:        name,
		phoneNumber: phoneNumber,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Getters
func (u *User) ID() string               { return u.id }
func (u *User) Email() Email             { return u.email }
func (u *User) Name() string             { return u.name }
func (u *User) PhoneNumber() *PhoneNumber { return u.phoneNumber }
func (u *User) CreatedAt() time.Time     { return u.createdAt }
func (u *User) UpdatedAt() time.Time     { return u.updatedAt }

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(name string, phoneNumber *PhoneNumber) {
	u.name = name
	u.phoneNumber = phoneNumber
	u.updatedAt = time.Now()
}

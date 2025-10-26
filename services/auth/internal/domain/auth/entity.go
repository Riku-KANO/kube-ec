package auth

import (
	"time"
)

// User represents an authenticated user entity
type User struct {
	id          string
	email       Email
	password    Password
	name        string
	phoneNumber string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewUser creates a new User entity
func NewUser(
	id string,
	email Email,
	password Password,
	name string,
	phoneNumber string,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:          id,
		email:       email,
		password:    password,
		name:        name,
		phoneNumber: phoneNumber,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// ID returns the user ID
func (u *User) ID() string {
	return u.id
}

// Email returns the user email
func (u *User) Email() Email {
	return u.email
}

// Password returns the user password
func (u *User) Password() Password {
	return u.password
}

// Name returns the user name
func (u *User) Name() string {
	return u.name
}

// PhoneNumber returns the user phone number
func (u *User) PhoneNumber() string {
	return u.phoneNumber
}

// CreatedAt returns the creation time
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the last update time
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// VerifyPassword verifies if the provided password matches
func (u *User) VerifyPassword(plainPassword string) error {
	return u.password.Verify(plainPassword)
}

// ChangePassword changes the user's password
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	// Verify old password
	if err := u.password.Verify(oldPassword); err != nil {
		return err
	}

	// Create new password
	newPwd, err := NewPassword(newPassword)
	if err != nil {
		return err
	}

	u.password = newPwd
	u.updatedAt = time.Now()
	return nil
}

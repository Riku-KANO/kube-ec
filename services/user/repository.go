package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	commonpb "github.com/Riku-KANO/kube-ec/pkg/proto/common"
	pb "github.com/Riku-KANO/kube-ec/services/user/proto"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

type UserRow struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	PhoneNumber  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash, name, phoneNumber string) (*UserRow, error) {
	query := `
		INSERT INTO users (id, email, password_hash, name, phone_number, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
		RETURNING id, email, password_hash, name, phone_number, created_at, updated_at
	`
	now := time.Now()
	user := &UserRow{}

	err := r.db.QueryRowContext(ctx, query,
		email,
		passwordHash,
		name,
		phoneNumber,
		now,
		now,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*UserRow, error) {
	query := `
		SELECT id, email, password_hash, name, phone_number, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &UserRow{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	return user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*UserRow, error) {
	query := `
		SELECT id, email, password_hash, name, phone_number, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &UserRow{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	return user, err
}

func (r *UserRepository) Update(ctx context.Context, id, name, phoneNumber string) error {
	query := `
		UPDATE users
		SET name = $2, phone_number = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		id,
		name,
		phoneNumber,
		time.Now(),
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func userRowToProto(user *UserRow) *pb.User {
	return &pb.User{
		Id:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		CreatedAt: &commonpb.Timestamp{
			Seconds: user.CreatedAt.Unix(),
		},
		UpdatedAt: &commonpb.Timestamp{
			Seconds: user.UpdatedAt.Unix(),
		},
	}
}

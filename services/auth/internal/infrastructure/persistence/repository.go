package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	pkgerrors "github.com/Riku-KANO/kube-ec/pkg/errors"
	"github.com/Riku-KANO/kube-ec/services/auth/internal/domain/auth"
	"github.com/lib/pq"
)

const (
	// dbTimeout is the default timeout for database operations
	dbTimeout = 10 * time.Second
)

// AuthRepository implements auth.Repository using PostgreSQL
type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new AuthRepository
func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *AuthRepository) CreateUser(ctx context.Context, user *auth.User) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
		INSERT INTO users (id, email, password_hash, name, phone_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID(),
		user.Email().String(),
		user.Password().Hash(),
		user.Name(),
		user.PhoneNumber(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	if err != nil {
		// Check for unique constraint violation (duplicate email)
		if isDuplicateKeyError(err) {
			return pkgerrors.ErrAlreadyExists
		}
		return pkgerrors.Wrap(pkgerrors.ErrInternal, fmt.Sprintf("failed to create user: %v", err))
	}

	return nil
}

// FindByEmail retrieves a user by email address
func (r *AuthRepository) FindByEmail(ctx context.Context, email auth.Email) (*auth.User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
		SELECT id, email, password_hash, name, phone_number, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var (
		id           string
		emailStr     string
		passwordHash string
		name         string
		phoneNumber  string
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&id,
		&emailStr,
		&passwordHash,
		&name,
		&phoneNumber,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkgerrors.ErrNotFound
	}
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInternal, fmt.Sprintf("failed to find user: %v", err))
	}

	return rowToUser(id, emailStr, passwordHash, name, phoneNumber, createdAt, updatedAt)
}

// FindByID retrieves a user by ID
func (r *AuthRepository) FindByID(ctx context.Context, id string) (*auth.User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
		SELECT id, email, password_hash, name, phone_number, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var (
		emailStr     string
		passwordHash string
		name         string
		phoneNumber  string
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&id,
		&emailStr,
		&passwordHash,
		&name,
		&phoneNumber,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkgerrors.ErrNotFound
	}
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInternal, fmt.Sprintf("failed to find user: %v", err))
	}

	return rowToUser(id, emailStr, passwordHash, name, phoneNumber, createdAt, updatedAt)
}

// UpdatePassword updates a user's password
func (r *AuthRepository) UpdatePassword(ctx context.Context, userID string, password auth.Password) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
		UPDATE users
		SET password_hash = $2, updated_at = $3
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID, password.Hash(), time.Now())
	if err != nil {
		return pkgerrors.Wrap(pkgerrors.ErrInternal, fmt.Sprintf("failed to update password: %v", err))
	}

	return nil
}

// rowToUser converts database row to domain User entity
func rowToUser(id, emailStr, passwordHash, name, phoneNumber string, createdAt, updatedAt time.Time) (*auth.User, error) {
	email, err := auth.NewEmail(emailStr)
	if err != nil {
		return nil, pkgerrors.Wrap(pkgerrors.ErrInternal, "invalid email in database")
	}

	password := auth.NewPasswordFromHash(passwordHash)

	return auth.NewUser(
		id,
		email,
		password,
		name,
		phoneNumber,
		createdAt,
		updatedAt,
	), nil
}

// isDuplicateKeyError checks if the error is a duplicate key constraint violation
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	// Check if error is a PostgreSQL error with code 23505 (unique_violation)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}

	return false
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/qper/hf/internal/domain"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, username, email, passwordHash string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, email
	`, username, email, passwordHash)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UserExists(ctx context.Context, username, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE (username = $1 OR email = $2) AND deleted_at IS NULL
		)
	`, username, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *AuthRepository) CreateRecoveryCodes(ctx context.Context, userID string, codeHashes []string) error {
	if len(codeHashes) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(codeHashes))
	args := make([]any, 0, len(codeHashes)*2+1)
	args = append(args, userID)
	for i, hash := range codeHashes {
		placeholders = append(placeholders, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, hash)
	}

	query := "INSERT INTO recovery_codes (user_id, code_hash) VALUES " + strings.Join(placeholders, ", ")
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

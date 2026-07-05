package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/qper/hf/internal/domain"
	"github.com/qper/hf/internal/service"
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

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash
		FROM users
		WHERE username = $1 AND deleted_at IS NULL
		LIMIT 1
	`, username)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
		LIMIT 1
	`, userID)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUnusedRecoveryCodes(ctx context.Context, userID string) ([]service.RecoveryCodeRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, code_hash
		FROM recovery_codes
		WHERE user_id = $1 AND used_at IS NULL
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var codes []service.RecoveryCodeRecord
	for rows.Next() {
		var code service.RecoveryCodeRecord
		if err := rows.Scan(&code.ID, &code.CodeHash); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return codes, nil
}

func (r *AuthRepository) MarkRecoveryCodeUsed(ctx context.Context, recoveryCodeID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE recovery_codes
		SET used_at = NOW()
		WHERE id = $1
	`, recoveryCodeID)
	return err
}

func (r *AuthRepository) DeleteRecoveryCodes(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM recovery_codes
		WHERE user_id = $1
	`, userID)
	return err
}

func (r *AuthRepository) CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO sessions (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (r *AuthRepository) GetSessionByToken(ctx context.Context, token string) (*service.SessionRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, token_hash, expires_at, revoked_at
		FROM sessions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var sessions []service.SessionRecord
	for rows.Next() {
		var session service.SessionRecord
		var revokedAt sql.NullTime
		if err := rows.Scan(&session.ID, &session.UserID, &session.TokenHash, &session.ExpiresAt, &revokedAt); err != nil {
			return nil, err
		}
		if revokedAt.Valid {
			session.RevokedAt = &revokedAt.Time
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, session := range sessions {
		if err := bcrypt.CompareHashAndPassword([]byte(session.TokenHash), []byte(token)); err == nil {
			return &session, nil
		}
	}
	return nil, nil
}

func (r *AuthRepository) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL
	`, sessionID)
	return err
}

func (r *AuthRepository) RevokeSessionByToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE token_hash = $1 AND revoked_at IS NULL
	`, token)
	return err
}

func (r *AuthRepository) RevokeAllSessions(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL
	`, userID)
	return err
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

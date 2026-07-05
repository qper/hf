package service

import (
	"context"
	"testing"
	"time"

	"github.com/qper/hf/internal/domain"
)

type fakeAuthRepository struct {
	secrets map[string]*SessionRecord
}

func (f *fakeAuthRepository) CreateUser(ctx context.Context, username, email, passwordHash string) (*domain.User, error) {
	return nil, nil
}
func (f *fakeAuthRepository) UserExists(ctx context.Context, username, email string) (bool, error) {
	return false, nil
}
func (f *fakeAuthRepository) CreateRecoveryCodes(ctx context.Context, userID string, codeHashes []string) error {
	return nil
}
func (f *fakeAuthRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return nil, nil
}
func (f *fakeAuthRepository) CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	return nil
}
func (f *fakeAuthRepository) GetSessionByToken(ctx context.Context, token string) (*SessionRecord, error) {
	return f.secrets[token], nil
}
func (f *fakeAuthRepository) RevokeSession(ctx context.Context, sessionID string) error {
	for _, session := range f.secrets {
		if session.ID == sessionID {
			now := time.Now().UTC()
			session.RevokedAt = &now
			break
		}
	}
	return nil
}
func (f *fakeAuthRepository) RevokeSessionByToken(ctx context.Context, token string) error {
	if session, ok := f.secrets[token]; ok {
		now := time.Now().UTC()
		session.RevokedAt = &now
	}
	return nil
}
func (f *fakeAuthRepository) RevokeAllSessions(ctx context.Context, userID string) error {
	for _, session := range f.secrets {
		if session.UserID == userID {
			now := time.Now().UTC()
			session.RevokedAt = &now
		}
	}
	return nil
}

func (f *fakeAuthRepository) DeleteRecoveryCodes(ctx context.Context, userID string) error {
	return nil
}

func (f *fakeAuthRepository) GetUnusedRecoveryCodes(ctx context.Context, userID string) ([]RecoveryCodeRecord, error) {
	return nil, nil
}

func (f *fakeAuthRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	return nil, nil
}

func (f *fakeAuthRepository) MarkRecoveryCodeUsed(ctx context.Context, recoveryCodeID string) error {
	return nil
}

func TestRefreshRejectsRevokedToken(t *testing.T) {
	ctx := context.Background()
	repo := &fakeAuthRepository{secrets: map[string]*SessionRecord{"old": {ID: "s1", UserID: "u1", TokenHash: "old", ExpiresAt: time.Now().UTC().Add(time.Hour)}}}
	service := &AuthService{repo: repo}
	_, err := service.Refresh(ctx, "old")
	if err == nil {
		t.Fatal("expected error for revoked token")
	}
}

package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"regexp"
	"strings"

	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/bcrypt"

	"github.com/qper/hf/internal/domain"
)

var (
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	ErrConflict     = fmt.Errorf("username or email already exists")
	ErrValidation   = fmt.Errorf("invalid registration payload")
)

type AuthService struct {
	repo AuthRepository
}

type AuthRepository interface {
	CreateUser(ctx context.Context, username, email, passwordHash string) (*domain.User, error)
	UserExists(ctx context.Context, username, email string) (bool, error)
	CreateRecoveryCodes(ctx context.Context, userID string, codeHashes []string) error
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, ErrValidation
	}

	exists, err := s.repo.UserExists(ctx, strings.TrimSpace(req.Username), strings.TrimSpace(req.Email))
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrConflict
	}

	passwordHash, err := argon2id.CreateHash(req.Password, &argon2id.Params{Memory: 64 * 1024, Iterations: 1, Parallelism: 2, SaltLength: 16, KeyLength: 32})
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, strings.TrimSpace(req.Username), strings.TrimSpace(req.Email), passwordHash)
	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, 8)
	codeHashes := make([]string, 0, 8)
	for i := 0; i < 8; i++ {
		code := generateRecoveryCode()
		codes = append(codes, code)
		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		codeHashes = append(codeHashes, string(hash))
	}

	if err := s.repo.CreateRecoveryCodes(ctx, user.ID, codeHashes); err != nil {
		return nil, err
	}

	return &domain.RegisterResponse{User: *user, RecoveryCodes: codes}, nil
}

func validateRegisterRequest(req domain.RegisterRequest) error {
	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(req.Email)
	password := req.Password

	if len(username) < 3 || len(username) > 50 || !usernamePattern.MatchString(username) {
		return fmt.Errorf("invalid username")
	}
	if !strings.Contains(email, "@") || len(strings.TrimSpace(email)) == 0 {
		return fmt.Errorf("invalid email")
	}
	if len(password) < 8 || !containsDigit(password) {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func containsDigit(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func generateRecoveryCode() string {
	b := make([]byte, 10)
	_, _ = rand.Read(b)
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	return strings.ToUpper(strings.TrimRight(encoded, "="))
}

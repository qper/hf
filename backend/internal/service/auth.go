package service

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base32"
	"encoding/pem"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/qper/hf/internal/domain"
)

var (
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	ErrConflict     = fmt.Errorf("username or email already exists")
	ErrValidation   = fmt.Errorf("invalid registration payload")
	ErrUnauthorized = fmt.Errorf("invalid credentials")
	ErrSession      = fmt.Errorf("invalid session")
)

type AuthService struct {
	repo AuthRepository
}

type AuthRepository interface {
	CreateUser(ctx context.Context, username, email, passwordHash string) (*domain.User, error)
	UserExists(ctx context.Context, username, email string) (bool, error)
	CreateRecoveryCodes(ctx context.Context, userID string, codeHashes []string) error
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	GetUnusedRecoveryCodes(ctx context.Context, userID string) ([]RecoveryCodeRecord, error)
	MarkRecoveryCodeUsed(ctx context.Context, recoveryCodeID string) error
	DeleteRecoveryCodes(ctx context.Context, userID string) error
	CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	GetSessionByToken(ctx context.Context, token string) (*SessionRecord, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeSessionByToken(ctx context.Context, token string) error
	RevokeAllSessions(ctx context.Context, userID string) error
}

type SessionRecord struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type RecoveryCodeRecord struct {
	ID       string
	CodeHash string
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

	codes, codeHashes, err := s.generateRecoveryCodes()
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateRecoveryCodes(ctx, user.ID, codeHashes); err != nil {
		return nil, err
	}

	return &domain.RegisterResponse{User: *user, RecoveryCodes: codes}, nil
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" || strings.TrimSpace(req.Password) == "" {
		return nil, ErrUnauthorized
	}

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	if _, err := argon2id.ComparePasswordAndHash(req.Password, user.PasswordHash); err != nil {
		return nil, ErrUnauthorized
	}

	accessToken, err := s.createAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken := newRefreshToken()
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	if err := s.repo.CreateSession(ctx, user.ID, string(hash), expiresAt); err != nil {
		return nil, err
	}

	return &domain.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "Bearer", ExpiresIn: 900}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*domain.RefreshResponse, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, ErrSession
	}

	session, err := s.repo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrSession
	}
	if session.RevokedAt != nil || session.ExpiresAt.Before(time.Now().UTC()) {
		if session.UserID != "" {
			_ = s.repo.RevokeAllSessions(ctx, session.UserID)
		}
		return nil, ErrSession
	}

	if err := s.repo.RevokeSession(ctx, session.ID); err != nil {
		return nil, err
	}

	newRefreshToken := newRefreshToken()
	hash, err := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	if err := s.repo.CreateSession(ctx, session.UserID, string(hash), expiresAt); err != nil {
		return nil, err
	}

	accessToken, err := s.createAccessToken(session.UserID)
	if err != nil {
		return nil, err
	}

	return &domain.RefreshResponse{AccessToken: accessToken, RefreshToken: newRefreshToken, TokenType: "Bearer", ExpiresIn: 900}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return ErrSession
	}
	return s.repo.RevokeSessionByToken(ctx, refreshToken)
}

func (s *AuthService) LogoutAll(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return ErrSession
	}
	session, err := s.repo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	if session == nil {
		return ErrSession
	}
	return s.repo.RevokeAllSessions(ctx, session.UserID)
}

func (s *AuthService) RecoverWithRecoveryCode(ctx context.Context, req domain.RecoverRequest) (*domain.RecoverResponse, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" || strings.TrimSpace(req.RecoveryCode) == "" {
		return nil, ErrUnauthorized
	}

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	unusedCodes, err := s.repo.GetUnusedRecoveryCodes(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	for _, code := range unusedCodes {
		if bcrypt.CompareHashAndPassword([]byte(code.CodeHash), []byte(req.RecoveryCode)) == nil {
			if err := s.repo.MarkRecoveryCodeUsed(ctx, code.ID); err != nil {
				return nil, err
			}

			accessToken, err := s.createAccessTokenWithMustChangePassword(user.ID)
			if err != nil {
				return nil, err
			}

			return &domain.RecoverResponse{AccessToken: accessToken, TokenType: "Bearer", ExpiresIn: 900}, nil
		}
	}

	return nil, ErrUnauthorized
}

func (s *AuthService) GetRecoveryCodeCount(ctx context.Context, userID string) (int, error) {
	codes, err := s.repo.GetUnusedRecoveryCodes(ctx, userID)
	if err != nil {
		return 0, err
	}
	return len(codes), nil
}

func (s *AuthService) RegenerateRecoveryCodes(ctx context.Context, userID, password string) (*domain.RecoveryCodeRegenerationResponse, error) {
	if strings.TrimSpace(password) == "" {
		return nil, ErrUnauthorized
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	if _, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash); err != nil {
		return nil, ErrUnauthorized
	}

	if err := s.repo.DeleteRecoveryCodes(ctx, userID); err != nil {
		return nil, err
	}

	codes, codeHashes, err := s.generateRecoveryCodes()
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateRecoveryCodes(ctx, userID, codeHashes); err != nil {
		return nil, err
	}

	return &domain.RecoveryCodeRegenerationResponse{RecoveryCodes: codes}, nil
}

func (s *AuthService) createAccessTokenWithMustChangePassword(userID string) (string, error) {
	privateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "secrets/jwt.key"
	}

	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("read private key: %w", err)
	}
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", fmt.Errorf("decode private key")
	}
	parsedKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse private key: %w", err)
	}

	claims := jwt.MapClaims{
		"sub":                  userID,
		"iat":                  time.Now().UTC().Unix(),
		"exp":                  time.Now().UTC().Add(15 * time.Minute).Unix(),
		"must_change_password": true,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(parsedKey)
}

func (s *AuthService) generateRecoveryCodes() ([]string, []string, error) {
	codes := make([]string, 0, 8)
	codeHashes := make([]string, 0, 8)
	for i := 0; i < 8; i++ {
		code := generateRecoveryCode()
		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			return nil, nil, err
		}
		codes = append(codes, code)
		codeHashes = append(codeHashes, string(hash))
	}
	return codes, codeHashes, nil
}

func (s *AuthService) createAccessToken(userID string) (string, error) {
	privateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "secrets/jwt.key"
	}

	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("read private key: %w", err)
	}
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", fmt.Errorf("decode private key")
	}
	parsedKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse private key: %w", err)
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(parsedKey)
}

func newRefreshToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
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

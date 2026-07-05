package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base32"
	"encoding/hex"
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
)

type AuthService struct {
	repo AuthRepository
}

type AuthRepository interface {
	CreateUser(ctx context.Context, username, email, passwordHash string) (*domain.User, error)
	UserExists(ctx context.Context, username, email string) (bool, error)
	CreateRecoveryCodes(ctx context.Context, userID string, codeHashes []string) error
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
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
	hash := sha256.Sum256([]byte(refreshToken))
	expiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	if err := s.repo.CreateSession(ctx, user.ID, hex.EncodeToString(hash[:]), expiresAt); err != nil {
		return nil, err
	}

	return &domain.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "Bearer", ExpiresIn: 900}, nil
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

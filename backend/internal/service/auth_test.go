package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateAccessTokenUsesPrivateKeyAndVerifiesWithPublicKey(t *testing.T) {
	tmpDir := t.TempDir()
	privatePath := filepath.Join(tmpDir, "jwt.key")
	publicPath := filepath.Join(tmpDir, "jwt.pub")

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	privateBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal private key: %v", err)
	}
	if err := os.WriteFile(privatePath, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateBytes}), 0o600); err != nil {
		t.Fatalf("write private key: %v", err)
	}
	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}
	if err := os.WriteFile(publicPath, pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicBytes}), 0o644); err != nil {
		t.Fatalf("write public key: %v", err)
	}

	oldPrivate := os.Getenv("JWT_PRIVATE_KEY_PATH")
	oldPublic := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if err := os.Setenv("JWT_PRIVATE_KEY_PATH", privatePath); err != nil {
		t.Fatalf("set private key path: %v", err)
	}
	if err := os.Setenv("JWT_PUBLIC_KEY_PATH", publicPath); err != nil {
		t.Fatalf("set public key path: %v", err)
	}
	defer func() {
		_ = os.Setenv("JWT_PRIVATE_KEY_PATH", oldPrivate)
		_ = os.Setenv("JWT_PUBLIC_KEY_PATH", oldPublic)
	}()

	svc := &AuthService{}
	token, err := svc.createAccessToken("user-1")
	if err != nil {
		t.Fatalf("create access token: %v", err)
	}

	publicPEM, err := os.ReadFile(publicPath)
	if err != nil {
		t.Fatalf("read public key: %v", err)
	}
	publicBlock, _ := pem.Decode(publicPEM)
	if publicBlock == nil {
		t.Fatalf("decode public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		t.Fatalf("parse public key: %v", err)
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return publicKey, nil
	})
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("unexpected claims type %T", parsed.Claims)
	}
	if claims["sub"] != "user-1" {
		t.Fatalf("expected subject user-1, got %v", claims["sub"])
	}
	if exp, ok := claims["exp"].(float64); !ok || exp <= float64(time.Now().Add(14*time.Minute).Unix()) {
		t.Fatalf("expected expiry in the future, got %v", claims["exp"])
	}
}

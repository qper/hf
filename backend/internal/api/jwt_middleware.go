package api

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var (
	pubKey     *ecdsa.PublicKey
	pubKeyOnce sync.Once
	pubKeyErr  error
)

const (
	ContextUserID             = "user_id"
	ContextMustChangePassword = "must_change_password"
)

func loadPublicKey() (*ecdsa.PublicKey, error) {
	pubKeyOnce.Do(func() {
		path := os.Getenv("JWT_PUBLIC_KEY_PATH")
		if path == "" {
			path = "secrets/jwt.pub"
		}
		b, err := os.ReadFile(path)
		if err != nil {
			pubKeyErr = err
			return
		}
		block, _ := pem.Decode(b)
		if block == nil {
			pubKeyErr = fmt.Errorf("invalid public key PEM")
			return
		}
		k, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			pubKeyErr = err
			return
		}
		pk, ok := k.(*ecdsa.PublicKey)
		if !ok {
			pubKeyErr = fmt.Errorf("public key is not ECDSA")
			return
		}
		pubKey = pk
	})
	return pubKey, pubKeyErr
}

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			pk, err := loadPublicKey()
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "public key unavailable"})
			}
			auth := c.Request().Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
			}
			tokenString := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
			parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
				if t.Method.Alg() != jwt.SigningMethodES256.Alg() {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
				}
				return pk, nil
			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			claims, ok := parsed.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			sub, ok := claims["sub"].(string)
			if !ok || strings.TrimSpace(sub) == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			if mustChange, ok := claims["must_change_password"].(bool); ok && mustChange {
				if c.Path() != "/api/v1/me/recovery-codes" {
					return c.JSON(http.StatusForbidden, map[string]string{"error": "password change required"})
				}
			}

			c.Set(ContextUserID, sub)
			return next(c)
		}
	}
}

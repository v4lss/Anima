// Package jwt wraps golang-jwt to sign and verify tokens for Animas.
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Service signs and verifies JWTs.
type Service struct {
	secret []byte
	expiry time.Duration
}

type claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

// New creates a JWT Service.
func New(secret string, expiry time.Duration) *Service {
	return &Service{secret: []byte(secret), expiry: expiry}
}

// Sign creates a signed token encoding the userID.
func (s *Service) Sign(userID string) (string, error) {
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(s.secret)
}

// Verify parses and validates the token, returning the embedded userID.
func (s *Service) Verify(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return "", err
	}
	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}
	return c.UserID, nil
}

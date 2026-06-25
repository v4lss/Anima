// Package auth (application layer) - use-case: Register.
package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/v4lss/animas/internal/domain/user"
)

// RegisterInput holds sign-up data from the client.
type RegisterInput struct {
	Email    string
	Password string
}

// Register creates a new user account.
func Register(ctx context.Context, repo user.Repository, in RegisterInput) (*user.User, error) {
	if strings.TrimSpace(in.Email) == "" || strings.TrimSpace(in.Password) == "" {
		return nil, errors.New("email and password are required")
	}

	existing, err := repo.FindByEmail(ctx, in.Email)
	if err == nil && existing != nil {
		return nil, errors.New("email already registered")
	}

	hash, err := user.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		Email:        in.Email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	if err := repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

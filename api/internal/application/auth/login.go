// Package auth (application layer) - use-case: Login.
// Verifies credentials and returns a signed JWT.
package auth

import (
	"context"
	"errors"

	"github.com/v4lss/animas/internal/domain/user"
	"github.com/v4lss/animas/pkg/jwt"
)

// LoginInput holds the credentials sent by the client.
type LoginInput struct {
	Email    string
	Password string
}

// LoginOutput is returned on success.
type LoginOutput struct {
	Token string
	User  *user.User
}

// Login authenticates a user and issues a JWT.
func Login(ctx context.Context, repo user.Repository, jwtSvc *jwt.Service, in LoginInput) (*LoginOutput, error) {
	u, err := repo.FindByEmail(ctx, in.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(in.Password, u.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := jwtSvc.Sign(u.ID)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{Token: token, User: u}, nil
}

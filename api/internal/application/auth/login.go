// Package auth (application layer) - use-case: Login.
// Verifies credentials and returns a signed JWT.
package auth

import (
	"context"
	"errors"
	"log"

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
	log.Printf("[Login] Attempting login for email: %s", in.Email)
	
	u, err := repo.FindByEmail(ctx, in.Email)
	if err != nil {
		log.Printf("[Login] User not found: %v", err)
		return nil, errors.New("invalid credentials")
	}

	log.Printf("[Login] User found, checking password")
	if !user.CheckPassword(in.Password, u.PasswordHash) {
		log.Printf("[Login] Password check failed")
		return nil, errors.New("invalid credentials")
	}

	token, err := jwtSvc.Sign(u.ID)
	if err != nil {
		log.Printf("[Login] JWT sign failed: %v", err)
		return nil, err
	}

	log.Printf("[Login] Login successful for user: %s", u.ID)
	return &LoginOutput{Token: token, User: u}, nil
}

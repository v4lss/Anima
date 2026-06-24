// Package user - repository port.
package user

import "context"

// Repository is the persistence interface for User.
type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

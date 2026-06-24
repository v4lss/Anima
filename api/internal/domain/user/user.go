// Package user defines the User entity.
package user

import "time"

// User represents a registered account.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

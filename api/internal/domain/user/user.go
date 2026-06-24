// Package user defines the User entity.
package user

import "time"

// User represents a registered account.
type User struct {
	ID           string    `bson:"id"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"passwordhash"`
	CreatedAt    time.Time `bson:"createdat"`
}

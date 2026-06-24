// Package validator provides basic request validation helpers.
package validator

import "regexp"

var emailRE = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// IsEmail returns true if s looks like a valid email address.
func IsEmail(s string) bool {
	return emailRE.MatchString(s)
}

// MinLength returns true if s has at least n runes.
func MinLength(s string, n int) bool {
	return len([]rune(s)) >= n
}

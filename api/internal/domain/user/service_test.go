package user

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Fatal("hash is empty")
	}

	if hash == password {
		t.Error("hash should not equal plain password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "mysecretpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for correct password")
	}

	if CheckPassword("wrongpassword", hash) {
		t.Error("CheckPassword should return false for incorrect password")
	}
}

func TestCheckPasswordInvalidHash(t *testing.T) {
	if CheckPassword("password", "invalidhash") {
		t.Error("CheckPassword should return false for invalid hash")
	}
}

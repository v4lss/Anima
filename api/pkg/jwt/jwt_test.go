package jwt

import (
	"testing"
	"time"
)

func TestSignAndVerify(t *testing.T) {
	svc := New("test-secret", 24*time.Hour)
	userID := "user123"

	token, err := svc.Sign(userID)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if token == "" {
		t.Fatal("token is empty")
	}

	decodedID, err := svc.Verify(token)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if decodedID != userID {
		t.Errorf("expected %s, got %s", userID, decodedID)
	}
}

func TestVerifyInvalidToken(t *testing.T) {
	svc := New("test-secret", 24*time.Hour)

	_, err := svc.Verify("invalid.token.here")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestVerifyWrongSecret(t *testing.T) {
	svc1 := New("secret1", 24*time.Hour)
	svc2 := New("secret2", 24*time.Hour)

	token, _ := svc1.Sign("user123")
	_, err := svc2.Verify(token)
	if err == nil {
		t.Error("expected error when verifying with wrong secret")
	}
}

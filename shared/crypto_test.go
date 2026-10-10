package shared

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	uuid, err := GenerateUUID()
	if err != nil {
		t.Fatalf("Failed to generate UUID: %v", err)
	}
	if len(uuid) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(uuid))
	}
}

func TestGenerateShortID(t *testing.T) {
	shortID, err := GenerateShortID()
	if err != nil {
		t.Fatalf("Failed to generate ShortID: %v", err)
	}
	if len(shortID) != 16 {
		t.Errorf("Expected ShortID length 16, got %d", len(shortID))
	}
}

func TestGeneratePassword(t *testing.T) {
	password, err := GeneratePassword()
	if err != nil {
		t.Fatalf("Failed to generate password: %v", err)
	}
	if len(password) != 32 {
		t.Errorf("Expected password length 32, got %d", len(password))
	}
	if _, err := hex.DecodeString(password); err != nil {
		t.Errorf("Password is not valid hex: %v", err)
	}
}

func TestGenerateX25519KeyPair(t *testing.T) {
	priv, pub, err := GenerateX25519KeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	privBytes, err := base64.RawURLEncoding.DecodeString(priv)
	if err != nil || len(privBytes) != 32 {
		t.Errorf("Invalid private key format/length")
	}

	pubBytes, err := base64.RawURLEncoding.DecodeString(pub)
	if err != nil || len(pubBytes) != 32 {
		t.Errorf("Invalid public key format/length")
	}
}

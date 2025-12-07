package pkg_test

import (
	"testing"

	"github.com/fazriegi/fintrack-be/internal/pkg"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "mysecurepassword"

	hash, err := pkg.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == password {
		t.Fatal("Hash should not equal the original password")
	}

	if !pkg.CheckPasswordHash(password, hash) {
		t.Fatal("CheckPasswordHash should return true for correct password")
	}

	if pkg.CheckPasswordHash("wrongpassword", hash) {
		t.Fatal("CheckPasswordHash should return false for incorrect password")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := "mysecretkey123"
	plain := "Hello World!"
	wrongKey := "wrong_key"

	encrypted, err := pkg.Encrypt(key, plain)
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}

	if encrypted == plain {
		t.Fatal("Encrypted text should not equal plain text")
	}

	decrypted, err := pkg.Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Decrypt returned error: %v", err)
	}

	if decrypted != plain {
		t.Errorf("Expected decrypted text to be %q, got %q", plain, decrypted)
	}

	_, err = pkg.Decrypt(wrongKey, encrypted)
	if err == nil {
		t.Fatal("Decrypt should fail when using the wrong key")
	}
}

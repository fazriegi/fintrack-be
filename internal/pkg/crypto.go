package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/fazriegi/fintrack-be/internal/infrastructure/config"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func deriveKey(keyStr string) []byte {
	h := sha256.Sum256([]byte(keyStr))
	return h[:] // 32 bytes for AES-256
}

func Encrypt(keyStr, value string) (string, error) {
	if keyStr == "" {
		keyStr = config.GetString("secret.encryptionKey")
		if keyStr == "" {
			return "", fmt.Errorf("encryption key is not set in config")
		}
	}

	key := deriveKey(keyStr)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", fmt.Errorf("invalid key size: must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	cipherData := aesGCM.Seal(nil, nonce, []byte(value), nil)
	output := append(nonce, cipherData...)

	return base64.StdEncoding.EncodeToString(output), nil
}

func Decrypt(keyStr, encodedCipher string) (string, error) {
	if keyStr == "" {
		keyStr = config.GetString("secret.encryptionKey")
		if keyStr == "" {
			return "", fmt.Errorf("encryption key is not set in config")
		}
	}

	key := deriveKey(keyStr)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", fmt.Errorf("invalid key size: must be 16, 24, or 32 bytes")
	}

	data, err := base64.StdEncoding.DecodeString(encodedCipher)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

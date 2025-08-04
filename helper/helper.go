package helper

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type ResponseJSON struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomState generates a cryptographically secure random state for OAuth
func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		slog.Error("Failed to generate random state", "error", err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateRandomPassword() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		slog.Error("Failed to generate random password", "error", err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

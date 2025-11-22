package helper

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"
)

type ResponseJSON struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginatedResponse struct {
	Items      any   `json:"items"`
	TotalCount int64 `json:"total_count"`
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

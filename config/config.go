package config

import (
	"fmt"
	"os"
)

func getSecretFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading secret file: %v", err)
	}
	return string(data), nil
}

func GetDBURL() (string, error) {
	filePath := os.Getenv("DB_URL_FILE")
	if filePath == "" {
		return "", fmt.Errorf("DB_URL_FILE environment variable is not set")
	}
	return getSecretFromFile(filePath)
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "5050"
	}
	return port
}

func GetJWTSecret() (string, error) {
	filePath := os.Getenv("JWT_SECRET_FILE")
	if filePath == "" {
		return "", fmt.Errorf("JWT_SECRET_FILE environment variable is not set")
	}
	return getSecretFromFile(filePath)
}

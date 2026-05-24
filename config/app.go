package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level         slog.Level
	Format        string
	SlowQueryTime int // Threshold for slow query logging in milliseconds
}

// Global logging configuration
var AppLogConfig LoggingConfig

// InitLogConfig initializes the logging configuration
func InitLogConfig() LoggingConfig {
	// Determine environment
	isDev := os.Getenv("GIN_MODE") != "release"

	// Configure log level from environment or use defaults
	logLevelStr := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	logLevel := slog.LevelInfo

	if isDev && logLevelStr == "" {
		logLevel = slog.LevelDebug
	} else {
		switch logLevelStr {
		case "DEBUG":
			logLevel = slog.LevelDebug
		case "INFO":
			logLevel = slog.LevelInfo
		case "WARN":
			logLevel = slog.LevelWarn
		case "ERROR":
			logLevel = slog.LevelError
		}
	}

	// Configure log format
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		if isDev {
			logFormat = "text"
		} else {
			logFormat = "json"
		}
	}

	// Configure slow query threshold
	slowQueryTime := 200 // Default 200ms
	if os.Getenv("SLOW_QUERY_TIME") != "" {
		_, _ = fmt.Sscanf(os.Getenv("SLOW_QUERY_TIME"), "%d", &slowQueryTime)
	}

	// Set global config
	AppLogConfig = LoggingConfig{
		Level:         logLevel,
		Format:        logFormat,
		SlowQueryTime: slowQueryTime,
	}

	return AppLogConfig
}

// getSecretFromFile reads and returns the content of a secret file as a string.
// It trims whitespace and ignores lines starting with '#' (comments).
func getSecretFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading secret file: %v", err)
	}

	// Trim whitespace
	content := strings.TrimSpace(string(data))

	// If the content starts with '#', it's a comment - return error to prompt user to update
	if strings.HasPrefix(content, "#") {
		return "", fmt.Errorf("connection string appears to be a comment/template. Please update secrets/db_url.txt")
	}

	return content, nil
}

// getSecretBytesFromFile reads and returns the content of a secret file as bytes
func getSecretBytesFromFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading secret file: %v", err)
	}
	return data, nil
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

func GetJWTSecret() ([]byte, error) {
	filePath := os.Getenv("JWT_SECRET_FILE")
	if filePath == "" {
		return nil, fmt.Errorf("JWT_SECRET_FILE environment variable is not set")
	}
	return getSecretBytesFromFile(filePath)
}

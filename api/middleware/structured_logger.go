package middleware

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// generateRequestID creates a simple request ID for tracing
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString generates a random string of specified length
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// StructuredLogger returns a Gin middleware that logs requests using structured logging.
// It includes:
// - Request ID for tracing
// - HTTP method and path
// - Status code
// - Response time in milliseconds
// - Client IP and User Agent
// - Any errors that occurred during handling
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID and add to context
		requestID := generateRequestID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration in milliseconds
		durationMs := float64(time.Since(start).Nanoseconds()) / 1e6

		// Get status code and determine log level
		status := c.Writer.Status()
		var logFn func(msg string, args ...any)

		switch {
		case status >= 500:
			logFn = slog.Error
		case status >= 400:
			logFn = slog.Warn
		default:
			logFn = slog.Info
		}

		// Build path with query parameters if present
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			path = path + "?" + raw
		}

		// Log request details
		logFn("http request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"duration_ms", durationMs,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"errors", c.Errors.String(),
		)
	}
}

package middleware

import (
	"crypto/rand"
	"strings"
	"time"

	"github.com/EdwinRincon/browersfc-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// RequestIDKey is the context key for the request ID
const RequestIDKey = "request_id"

// generateRequestID creates a simple request ID for tracing
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString generates a cryptographically secure random string of specified length
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	randBytes := make([]byte, n)
	if _, err := rand.Read(randBytes); err != nil {
		// Fallback to timestamp-based generation if crypto/rand fails
		return time.Now().Format("150405")
	}
	for i, rb := range randBytes {
		b[i] = letters[int(rb)%len(letters)]
	}
	return string(b)
}

// shouldSkipPath determines if a path should have minimal logging
func shouldSkipPath(path string) bool {
	skipPaths := []string{
		"/healthz",
		"/metrics",
		"/favicon.ico",
		"/css/",
		"/js/",
		"/static/",
	}

	for _, prefix := range skipPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// StructuredLogger returns a Gin middleware that logs requests using structured logging.
// It includes:
// - Request ID for tracing (respects incoming X-Request-ID header)
// - HTTP method and path
// - Status code
// - Response time in milliseconds
// - Client IP and User Agent
// - Any errors that occurred during handling
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for existing request ID from upstream services (e.g., Nginx)
		// or generate a new one if not present
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Store in context and header for downstream services
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration in milliseconds
		durationMs := float64(time.Since(start).Nanoseconds()) / 1e6

		// Use the centralized logger to log this HTTP request
		logger.LogHTTPRequest(c, durationMs, shouldSkipPath(c.Request.URL.Path))
	}
}

package middleware

import (
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/EdwinRincon/browersfc-api/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Context keys for storing user information in Gin context
const (
	usernameKey  = "username"
	roleKey      = "role"
	bearerPrefix = "Bearer "
)

// JwtAuthMiddleware authenticates requests using JWT tokens through the domain service.
// It supports both cookie-based authentication (preferred) and Authorization header for API calls.
func JwtAuthMiddleware(authService *service.AuthenticationDomainService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// First, try to get token from cookie (preferred method)
		if cookie, err := c.Cookie("token"); err == nil && cookie != "" {
			tokenString = cookie
		} else {
			// Fallback to Authorization header for API calls
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				helper.WriteErrorResponse(c, helper.NewUnauthorizedError("Authentication required"))
				c.Abort()
				return
			}

			if !strings.HasPrefix(authHeader, bearerPrefix) {
				helper.WriteErrorResponse(c, helper.NewUnauthorizedError("Invalid authorization header format"))
				c.Abort()
				return
			}

			// Extract token without the "Bearer " prefix
			tokenString = strings.TrimPrefix(authHeader, bearerPrefix)
		}

		// Validate token using domain service
		authClaims, err := authService.ValidateAuthentication(c.Request.Context(), tokenString)
		if err != nil {
			// Log detailed error for debugging but return generic message to client
			logger.Debug(c, "JWT token validation failed",
				"error", err.Error(),
				"ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent())

			helper.WriteErrorResponse(c, helper.NewUnauthorizedError("Invalid or expired token"))
			c.Abort()
			return
		}

		// Store claims in context using string keys
		c.Set(usernameKey, authClaims.Username)
		c.Set(roleKey, authClaims.Role)

		logger.Debug(c, "Successfully authenticated request",
			"username", authClaims.Username,
			"role", authClaims.Role)

		c.Next()
	}
}

// RBACMiddleware implements Role-Based Access Control by checking if the
// authenticated user's role is allowed to access the requested resource.
func RBACMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(roleKey)
		if !exists {
			helper.WriteErrorResponse(c, helper.NewUnauthorizedError("Authentication required"))
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(errors.New("invalid role format")))
			c.Abort()
			return
		}

		username, _ := c.Get(usernameKey)

		// Check if user's role is allowed
		if slices.Contains(allowedRoles, userRole) {
			logger.Info(c, "Access granted",
				"username", username,
				"role", userRole,
				"path", c.Request.URL.Path,
				"method", c.Request.Method)
			c.Next()
			return
		}

		// Log access denied event as warning since it may indicate security issues
		logger.Warn(c, "Access denied",
			"username", username,
			"role", userRole,
			"allowedRoles", allowedRoles,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"ip", c.ClientIP())

		helper.WriteErrorResponse(c, helper.NewForbiddenError("Insufficient permissions"))
		c.Abort()
	}
}

// SecurityHeadersMiddleware returns a middleware function that sets various security headers.
//
// X-Frame-Options: Prevents clickjacking attacks.
//
// Content-Security-Policy: Helps prevent cross-site scripting (XSS) and data injection attacks.
//
// X-XSS-Protection: Enables the browser's built-in XSS filter.
//
// Strict-Transport-Security: Enforces secure (HTTPS) connections to the server.
//
// Referrer-Policy: Controls the amount of referrer information sent along with requests.
//
// X-Content-Type-Options: Prevents MIME-sniffing a response away from the declared content-type.
//
// Permissions-Policy: Allows you to enable or disable various browser features and APIs.
// SecurityHeadersMiddleware adds security-related HTTP headers to all responses
// to protect against various common web vulnerabilities.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	// Pre-compute CSP header since it rarely changes
	cspHeader := config.BuildCSPHeader()

	return func(c *gin.Context) {
		// Host check in development only - allow any localhost port
		if config.Config.IsDevelopment {
			if !strings.HasPrefix(c.Request.Host, "localhost") {
				logger.Warn(c, "Invalid host header in development",
					"host", c.Request.Host,
					"ip", c.ClientIP())
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host"})
				return
			}
		}

		// Security Headers
		headers := map[string]string{
			"X-Frame-Options":                   "DENY",
			"Content-Security-Policy":           cspHeader,
			"Strict-Transport-Security":         "max-age=31536000; includeSubDomains; preload",
			"Referrer-Policy":                   "strict-origin-when-cross-origin",
			"X-Content-Type-Options":            "nosniff",
			"X-DNS-Prefetch-Control":            "off",
			"X-Download-Options":                "noopen",
			"X-Permitted-Cross-Domain-Policies": "none",
			"Permissions-Policy":                "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=(), interest-cohort=()",
			"Cross-Origin-Opener-Policy":        "same-origin",
			"Cross-Origin-Embedder-Policy":      "require-corp",
			"Cross-Origin-Resource-Policy":      "same-origin",
			"Expect-CT":                         "max-age=86400, enforce",
		}

		// Set all headers at once
		for key, value := range headers {
			c.Header(key, value)
		}

		c.Next()
	}
}

// CORSMiddleware returns a middleware to handle Cross-Origin Resource Sharing (CORS).
// It uses configuration from the app config to determine allowed origins for different environments.
func CORSMiddleware() gin.HandlerFunc {
	var allowedOrigins []string
	if config.Config.IsDevelopment {
		allowedOrigins = []string{
			"http://localhost:4200", // Angular dev server
		}
	} else {
		// Production origins
		allowedOrigins = []string{
			"https://app.browersfc.com",
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    false, // Explicitly disable wildcard matching for security
	})
}

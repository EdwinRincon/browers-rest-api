package middleware

import (
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/EdwinRincon/browersfc-api/pkg/jwt"
	"github.com/EdwinRincon/browersfc-api/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// contextKey type to avoid collisions in context values
type contextKey string

const (
	usernameKey contextKey = "username"
	roleKey     contextKey = "role"
)

// JwtAuthMiddleware authenticates requests using JWT tokens.
// It extracts the token from the Authorization header, validates it,
// and sets the authenticated user's information in the context.
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helper.RespondWithError(c, helper.Unauthorized("Authorization header is required"))
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			helper.RespondWithError(c, helper.Unauthorized("Invalid authorization header format"))
			c.Abort()
			return
		}

		// Extract token without the "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Get JWT secret and create validator
		jwtSecret, err := config.GetJWTSecret()
		if err != nil {
			helper.RespondWithError(c, helper.InternalError(err))
			c.Abort()
			return
		}

		validator := jwt.NewTokenValidator(jwtSecret)
		claims, err := validator.ValidateToken(tokenString)
		if err != nil {
			helper.RespondWithError(c, helper.Unauthorized(err.Error()))
			c.Abort()
			return
		}

		// Store claims in context using type-safe keys
		c.Set(string(usernameKey), claims.Username)
		c.Set(string(roleKey), claims.Role)

		logger.Debug(c, "Successfully authenticated request",
			"username", claims.Username,
			"role", claims.Role)

		c.Next()
	}
}

// RBACMiddleware implements Role-Based Access Control by checking if the
// authenticated user's role is allowed to access the requested resource.
func RBACMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(string(roleKey))
		if !exists {
			helper.RespondWithError(c, helper.Unauthorized("Authentication required"))
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			helper.RespondWithError(c, helper.InternalError(errors.New("invalid role format")))
			c.Abort()
			return
		}

		username, _ := c.Get(string(usernameKey))

		// Check if user's role is allowed
		if slices.Contains(allowedRoles, userRole) {
			logger.Debug(c, "Access granted",
				"username", username,
				"role", userRole,
				"path", c.Request.URL.Path,
				"method", c.Request.Method)
			c.Next()
			return
		}

		// Log access denied event
		logger.Info(c, "Access denied",
			"username", username,
			"role", userRole,
			"allowedRoles", allowedRoles,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"ip", c.ClientIP())

		helper.RespondWithError(c, helper.StatusForbidden("Insufficient permissions"))
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
		// Host check in development only
		if config.Config.IsDevelopment {
			if c.Request.Host != "localhost:5050" {
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

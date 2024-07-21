package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
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
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		expectedHost := "localhost:5050"
		if c.Request.Host != expectedHost {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		//TODO: Review the security headers
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	}
}

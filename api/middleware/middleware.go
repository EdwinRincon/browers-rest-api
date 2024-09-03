package middleware

import (
	"net/http"
	"strings"

	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/helper"
	jwtClaims "github.com/EdwinRincon/browersfc-api/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JwtAuthMiddleware es un middleware que verifica el token JWT
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helper.HandleError(c, helper.NewAppError(http.StatusUnauthorized, "Authorization header is required", ""), false)
			c.Abort()
			return
		}

		// Extraer el token del encabezado Authorization
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Si no se elimina "Bearer ", tokenString será igual a authHeader
			helper.HandleError(c, helper.NewAppError(http.StatusUnauthorized, "Invalid authorization header format", ""), false)
			c.Abort()
			return
		}

		// Obtener el secreto JWT
		jwtSecret, err := config.GetJWTSecret()
		if err != nil {
			helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to read JWT secret", err.Error()), false)
			c.Abort()
			return
		}

		// Parsear el token
		token, err := jwt.ParseWithClaims(tokenString, &jwtClaims.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			helper.HandleError(c, helper.NewAppError(http.StatusUnauthorized, "Invalid token", err.Error()), false)
			c.Abort()
			return
		}

		// Validar y extraer los claims
		claims, ok := token.Claims.(*jwtClaims.AppClaims)
		if !ok {
			helper.HandleError(c, helper.NewAppError(http.StatusUnauthorized, "Invalid token claims", ""), false)
			c.Abort()
			return
		}

		// Establecer los claims en el contexto
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

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

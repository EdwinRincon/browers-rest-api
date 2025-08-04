package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AppClaims defines JWT token claims structure
type AppClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token generation and validation
type JWTService struct {
	SecretKey []byte
}

// NewJWTService creates a new instance of JWTService
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{SecretKey: []byte(secretKey)}
}

// GenerateToken creates a new JWT token for a given user
func (s *JWTService) GenerateToken(username, role string) (string, error) {
	claims := &AppClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Token expires in 1 hour
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.SecretKey)
}

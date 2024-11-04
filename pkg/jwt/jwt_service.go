package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AppClaims define las reclamaciones del JWT
type AppClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTService maneja la generación y validación de JWTs
type JWTService struct {
	SecretKey []byte
}

// NewJWTService crea una nueva instancia de JWTService
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{SecretKey: []byte(secretKey)}
}

// GenerateToken genera un nuevo token JWT para un usuario dado
func (s *JWTService) GenerateToken(username, role string) (string, error) {
	claims := &AppClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Token válido por 1 hora
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.SecretKey)
}

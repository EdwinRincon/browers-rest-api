package jwt

import (
	"log"
	"time"

	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken genera un nuevo token JWT para un usuario dado
func GenerateToken(username, role string) (string, error) {
	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		log.Fatalf("Failed to read JWT secret from file: %v", err)
	}
	claims := &AppClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)), // Token válido por 2 horas
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

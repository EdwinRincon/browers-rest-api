package jwt

import "github.com/golang-jwt/jwt/v5"

type AppClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

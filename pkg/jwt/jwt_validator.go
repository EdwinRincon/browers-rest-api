package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenRequired  = errors.New("authorization token is required")
	ErrInvalidFormat  = errors.New("invalid authorization header format")
	ErrInvalidToken   = errors.New("invalid token")
	ErrInvalidClaims  = errors.New("invalid token claims")
	ErrUnsupportedAlg = errors.New("unsupported signing algorithm")
)

type TokenValidator struct {
	jwtSecret []byte
}

func NewTokenValidator(secret []byte) *TokenValidator {
	return &TokenValidator{
		jwtSecret: secret,
	}
}

func (v *TokenValidator) ValidateToken(tokenString string) (*AppClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AppClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnsupportedAlg
			}
			return v.jwtSecret, nil
		},
	)

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AppClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

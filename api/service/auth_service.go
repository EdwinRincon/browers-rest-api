package service

import (
	"github.com/EdwinRincon/browersfc-api/api/repository"
	"github.com/EdwinRincon/browersfc-api/pkg/jwt"
)

// AuthService defines the authentication service interface
type AuthService interface {
	GenerateToken(username, role string) (string, error)
}

// authService implements the AuthService interface
type authService struct {
	UserRepository repository.UserRepository
	JWTService     *jwt.JWTService
}

// NewAuthService creates a new authentication service instance
func NewAuthService(userRepo repository.UserRepository, jwtService *jwt.JWTService) AuthService {
	return &authService{
		UserRepository: userRepo,
		JWTService:     jwtService,
	}
}

// GenerateToken creates a new JWT token for an authenticated user
func (s *authService) GenerateToken(username, role string) (string, error) {
	return s.JWTService.GenerateToken(username, role)
}

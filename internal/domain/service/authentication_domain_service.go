package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// AuthenticationDomainService contains the business logic for authentication operations.
// It operates on domain entities and implements authentication business rules.
type AuthenticationDomainService struct {
	authRepository domain.AuthenticationRepository
}

// NewAuthenticationDomainService creates a new AuthenticationDomainService.
func NewAuthenticationDomainService(authRepository domain.AuthenticationRepository) *AuthenticationDomainService {
	return &AuthenticationDomainService{
		authRepository: authRepository,
	}
}

// GenerateToken creates a new authentication token for the given user.
func (s *AuthenticationDomainService) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	if user == nil {
		return "", constants.ErrInvalidData
	}

	return s.authRepository.GenerateAccessToken(ctx, user)
}

// ValidateAuthentication validates a token and returns authentication claims.
func (s *AuthenticationDomainService) ValidateAuthentication(ctx context.Context, token string) (*domain.AuthenticationClaims, error) {
	if token == "" {
		return nil, constants.ErrInvalidData
	}

	return s.authRepository.ValidateAccessToken(ctx, token)
}

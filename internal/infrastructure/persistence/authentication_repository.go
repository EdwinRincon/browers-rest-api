package persistence

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/adapter/mapper"
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/config"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/pkg/jwt"
)

// AuthenticationRepository implements domain.AuthenticationRepository interface.
// It handles the infrastructure concerns for authentication operations.
type AuthenticationRepository struct {
	jwtService     *jwt.JWTService
	mapper         *mapper.AuthenticationMapper
	roleRepository domain.RoleRepository
}

// NewAuthenticationRepository creates a new AuthenticationRepository instance.
func NewAuthenticationRepository(roleRepository domain.RoleRepository) *AuthenticationRepository {
	// Get JWT secret from configuration
	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		panic("Failed to get JWT secret: " + err.Error())
	}

	return &AuthenticationRepository{
		jwtService:     jwt.NewJWTService(string(jwtSecret)),
		mapper:         mapper.NewAuthenticationMapper(),
		roleRepository: roleRepository,
	}
}

// GenerateAccessToken creates a new JWT access token for the given user.
func (r *AuthenticationRepository) GenerateAccessToken(ctx context.Context, user *domain.User) (string, error) {
	if user == nil {
		return "", constants.ErrInvalidData
	}

	// Fetch the actual role name from the database using RoleID
	role, err := r.roleRepository.GetRoleByID(ctx, user.RoleID)
	if err != nil {
		return "", err
	}

	if role == nil {
		return "", constants.ErrRecordNotFound
	}

	// Generate token using the actual role name from database
	return r.jwtService.GenerateToken(user.Username, role.Name)
}

// ValidateAccessToken validates a JWT token and returns authentication claims.
func (r *AuthenticationRepository) ValidateAccessToken(ctx context.Context, token string) (*domain.AuthenticationClaims, error) {
	if token == "" {
		return nil, constants.ErrInvalidData
	}

	// Create validator with JWT secret
	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		return nil, err
	}

	validator := jwt.NewTokenValidator(jwtSecret)
	appClaims, err := validator.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Convert JWT claims to domain authentication claims
	return r.mapper.ToDomainClaims(appClaims), nil
}

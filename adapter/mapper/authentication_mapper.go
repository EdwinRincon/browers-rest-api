package mapper

import (
	"time"

	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/pkg/jwt"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

// AuthenticationMapper handles mapping between JWT claims and domain authentication
type AuthenticationMapper struct{}

// NewAuthenticationMapper creates a new instance of AuthenticationMapper
func NewAuthenticationMapper() *AuthenticationMapper {
	return &AuthenticationMapper{}
}

// ToDomainClaims converts JWT AppClaims to domain authentication claims
func (m *AuthenticationMapper) ToDomainClaims(appClaims *jwt.AppClaims) *domain.AuthenticationClaims {
	var expiresAt time.Time
	if appClaims.ExpiresAt != nil {
		expiresAt = appClaims.ExpiresAt.Time
	}

	var issuedAt time.Time
	if appClaims.IssuedAt != nil {
		issuedAt = appClaims.IssuedAt.Time
	} else {
		issuedAt = time.Now()
	}

	return &domain.AuthenticationClaims{
		UserID:    appClaims.Subject, // JWT Subject contains the user ID
		Username:  appClaims.Username,
		Role:      appClaims.Role,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}
}

// FromDomainClaims converts domain authentication claims to JWT AppClaims
func (m *AuthenticationMapper) FromDomainClaims(domainClaims *domain.AuthenticationClaims) *jwt.AppClaims {
	return &jwt.AppClaims{
		Username: domainClaims.Username,
		Role:     domainClaims.Role,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Subject:   domainClaims.UserID,
			ExpiresAt: jwtlib.NewNumericDate(domainClaims.ExpiresAt),
			IssuedAt:  jwtlib.NewNumericDate(domainClaims.IssuedAt),
		},
	}
}

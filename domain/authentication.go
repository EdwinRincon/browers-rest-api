package domain

import (
	"time"
)

// AuthenticationClaims represents the authentication information in the domain layer.
// It encapsulates the user's identity and permissions
type AuthenticationClaims struct {
	UserID    string
	Username  string
	Role      string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// IsValid checks if the authentication claims are valid and not expired.
func (a *AuthenticationClaims) IsValid() bool {
	return a.UserID != "" &&
		a.Username != "" &&
		!a.IssuedAt.IsZero() &&
		!a.ExpiresAt.IsZero() &&
		time.Now().Before(a.ExpiresAt)
}

// IsExpired checks if the authentication token has expired.
func (a *AuthenticationClaims) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}

// TimeUntilExpiry returns the duration until the token expires.
func (a *AuthenticationClaims) TimeUntilExpiry() time.Duration {
	if a.IsExpired() {
		return 0
	}
	return time.Until(a.ExpiresAt)
}

// HasRole checks if the authenticated user has the specified role.
func (a *AuthenticationClaims) HasRole(role string) bool {
	return a.Role == role
}

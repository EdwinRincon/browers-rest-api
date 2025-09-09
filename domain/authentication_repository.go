package domain

import (
	"context"
)

// AuthenticationRepository defines the interface for authentication operations.
// This port belongs in the domain layer
type AuthenticationRepository interface {
	// GenerateAccessToken creates a new authentication token for the given user
	GenerateAccessToken(ctx context.Context, user *User) (string, error)

	// ValidateAccessToken validates a token and returns the authentication claims
	ValidateAccessToken(ctx context.Context, token string) (*AuthenticationClaims, error)
}

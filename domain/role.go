package domain

import (
	"strings"
	"time"
)

// Role represents a user role in the domain layer.
type Role struct {
	ID          uint64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Common role names
const (
	RoleAdmin  = "admin"
	RolePlayer = "player"
	RoleCoach  = "coach"
)

// IsSystemRole returns true if this is a built-in system role.
func (r *Role) IsSystemRole() bool {
	roleName := strings.ToLower(r.Name)
	return roleName == RoleAdmin || roleName == RolePlayer || roleName == RoleCoach
}

// IsValid performs basic domain validation for the role.
func (r *Role) IsValid() bool {
	return r.Name != "" &&
		len(r.Name) <= 20 &&
		len(r.Description) <= 100 &&
		r.isValidRoleName()
}

// isValidRoleName checks if the role name contains only valid characters.
func (r *Role) isValidRoleName() bool {
	if r.Name == "" {
		return false
	}

	// Role names should contain only letters, numbers, and underscores
	for _, char := range r.Name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return false
		}
	}

	return true
}

// GetDisplayName returns a human-readable version of the role name.
func (r *Role) GetDisplayName() string {
	if r.Description != "" {
		return r.Description
	}

	// Capitalize first letter of role name
	if len(r.Name) > 0 {
		return strings.ToUpper(string(r.Name[0])) + strings.ToLower(r.Name[1:])
	}

	return r.Name
}

// CanBeDeleted returns true if the role can be safely deleted.
func (r *Role) CanBeDeleted() bool {
	// System roles cannot be deleted
	return !r.IsSystemRole()
}

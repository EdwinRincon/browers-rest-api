package constants

import (
	"errors"
)

// User Roles
const (
	RoleAdmin   = "admin"
	RolePlayer  = "player"
	RoleDefault = "fan"
)

// Common error messages
const (
	MsgInvalidID         = "Invalid ID format"
	MsgInvalidTeamID     = "Invalid team ID"
	MsgInvalidPlayerID   = "Invalid player ID"
	MsgInvalidUserID     = "Invalid user ID"
	MsgInvalidRoleData   = "Invalid role data"
	MsgInvalidSeasonData = "Invalid season data"
	MsgInvalidData       = "Invalid data"
	MsgInvalidTeamData   = "Invalid team data"
	MsgInvalidUserData   = "Invalid user data"
	MsgNotFound          = "Resource not found"
	MsgUnauthorized      = "Unauthorized access"
	MsgForbidden         = "Forbidden access"
	MsgInternalError     = "An unexpected error occurred"
)

// Error messages
var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrRecordAlreadyExists = errors.New("record already exists")
)

const APIBasePath = "/api"

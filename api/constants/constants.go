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
	ErrPlayerNotFound      = errors.New("player not found")
	ErrTeamNotFound        = errors.New("team not found")
	ErrSeasonNotFound      = errors.New("season not found")
	ErrOverlappingDates    = errors.New("date range overlaps with existing player team record")
)

const APIBasePath = "/api"

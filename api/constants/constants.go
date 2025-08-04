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

// Pagination
const (
	DefaultPageSize = 10
)

// Error Variables
var (
	ErrUserNotFound           = errors.New("user not found")
	ErrArticleNotFound        = errors.New("article not found")
	ErrPlayerNotFound         = errors.New("player not found")
	ErrRoleNotFound           = errors.New("role not found")
	ErrTeamNotFound           = errors.New("team not found")
	ErrMatchNotFound          = errors.New("match not found")
	ErrSeasonNotFound         = errors.New("season not found")
	ErrClassificationNotFound = errors.New("classification not found")
	ErrLineupNotFound         = errors.New("lineup not found")
	ErrTeamStatsNotFound      = errors.New("team stats not found")

	ErrInvalidUUID        = errors.New("invalid UUID")
	ErrTimezoneLoad       = errors.New("failed to load location")
	ErrUserUpdate         = errors.New("failed to update user")
	ErrUserDelete         = errors.New("failed to delete user")
	ErrCreateUser         = errors.New("failed to create user")
	ErrDuplicatedUsername = errors.New("duplicated username")
)

// Common Error Messages
const (
	ErrInvalidInput     = "Invalid input"
	ErrInvalidPlayerID  = "Invalid player ID"
	ErrInvalidArticleID = "Invalid article ID"
	ErrInvalidMatchID   = "Invalid match ID"
	ErrInvalidTeamID    = "Invalid team ID"
	ErrInvalidSeasonID  = "Invalid season ID"
)

const APIBasePath = "/api"

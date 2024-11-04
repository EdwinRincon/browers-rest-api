package constants

import "errors"

// User Roles
const (
	RoleAdmin  = "Admin"
	RolePlayer = "Player"
	RoleCoach  = "Coach"
)

// Pagination
const (
	DefaultPageSize = 10
)

// Error Variables
var ErrUserNotFound = errors.New("user not found")
var ErrArticleNotFound = errors.New("article not found")
var ErrPlayerNotFound = errors.New("player not found")
var ErrRoleNotFound = errors.New("role not found")
var ErrTeamNotFound = errors.New("team not found")
var ErrMatchNotFound = errors.New("match not found")
var ErrSeasonNotFound = errors.New("season not found")
var ErrClassificationNotFound = errors.New("classification not found")

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

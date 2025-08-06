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

// Error messages
var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrRecordAlreadyExists = errors.New("record already exists")
)

const APIBasePath = "/api"

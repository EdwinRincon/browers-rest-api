package domain

import (
	"context"
	"time"
)

// OverlapCheckData holds the data needed for date overlap validation checks.
// This struct is used in business logic validation.
type OverlapCheckData struct {
	PlayerID  uint64
	TeamID    uint64
	SeasonID  uint64
	StartDate time.Time
	EndDate   *time.Time
	IsUpdate  bool
	ID        uint64
}

// PlayerTeamRepository defines the interface for player-team association persistence operations.
// This repository interface belongs in the domain layer following hexagonal architecture principles.
type PlayerTeamRepository interface {
	Create(ctx context.Context, playerTeam *PlayerTeam) error
	GetPlayerTeamByID(ctx context.Context, id uint64) (*PlayerTeam, error)
	UpdatePlayerTeam(ctx context.Context, playerTeam *PlayerTeam) error
	DeletePlayerTeam(ctx context.Context, id uint64) error
	GetPaginatedPlayerTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]PlayerTeam, int64, error)

	// Relationship-specific query methods
	GetByPlayerID(ctx context.Context, playerID uint64) ([]PlayerTeam, error)
	GetPlayerTeamsByTeamID(ctx context.Context, teamID uint64) ([]PlayerTeam, error)
	GetPlayerTeamsBySeasonID(ctx context.Context, seasonID uint64) ([]PlayerTeam, error)

	DeleteByPlayerID(ctx context.Context, playerID uint64) error
	CheckOverlappingDates(ctx context.Context, data OverlapCheckData) (bool, error)
}

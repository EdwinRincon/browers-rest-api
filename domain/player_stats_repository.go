package domain

import "context"

// PlayerStatsRepository defines the interface for player statistics persistence operations.
// This port belongs in the domain layer
type PlayerStatsRepository interface {
	CreatePlayerStat(ctx context.Context, playerStat *PlayerStat) error
	GetPlayerStatByID(ctx context.Context, id uint64) (*PlayerStat, error)
	GetPlayerStatsByPlayerID(ctx context.Context, playerID uint64) ([]PlayerStat, error)
	GetPlayerStatsByMatchID(ctx context.Context, matchID uint64) ([]PlayerStat, error)
	GetPlayerStatsBySeasonID(ctx context.Context, seasonID uint64) ([]PlayerStat, error)
	GetPaginatedPlayerStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]PlayerStat, int64, error)
	UpdatePlayerStat(ctx context.Context, id uint64, playerStat *PlayerStat) error
	DeletePlayerStat(ctx context.Context, id uint64) error
}

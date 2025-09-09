package domain

import (
	"context"
)

// TeamStatsRepository defines the interface for team stats persistence operations.
// This port belongs in the domain layer
type TeamStatsRepository interface {
	CreateTeamStats(ctx context.Context, teamStats *TeamStats) error
	GetTeamStatsByID(ctx context.Context, id uint64) (*TeamStats, error)
	GetTeamStatsBySeasonAndTeam(ctx context.Context, seasonID, teamID uint64) (*TeamStats, error)
	GetTeamStatsBySeasonID(ctx context.Context, seasonID uint64) ([]TeamStats, error)
	GetTeamStatsByTeamID(ctx context.Context, teamID uint64) ([]TeamStats, error)
	GetPaginatedTeamStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]TeamStats, int64, error)
	UpdateTeamStats(ctx context.Context, id uint64, teamStats *TeamStats) error
	DeleteTeamStats(ctx context.Context, id uint64) error
}

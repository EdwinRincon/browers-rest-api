package domain

import (
	"context"
)

// MatchRepository defines the interface for match persistence operations.
// This port belongs in the domain layer
type MatchRepository interface {
	CreateMatch(ctx context.Context, match *Match) error
	GetMatchByID(ctx context.Context, id uint64) (*Match, error)
	GetPaginatedMatches(ctx context.Context, sort string, order string, page int, pageSize int) ([]Match, int64, error)
	GetMatchesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]Match, int64, error)
	GetMatchesByTeamID(ctx context.Context, teamID uint64, sort string, order string, page int, pageSize int) ([]Match, int64, error)
	GetNextMatchByTeamID(ctx context.Context, teamID uint64) (*Match, error)
	GetDetailedMatchByID(ctx context.Context, id uint64) (*Match, error)
	UpdateMatch(ctx context.Context, id uint64, match *Match) error
	DeleteMatch(ctx context.Context, id uint64) error
}

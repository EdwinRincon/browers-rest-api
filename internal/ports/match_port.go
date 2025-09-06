package ports

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
)

// MatchPort defines the interface for match data access operations.
// This port separates the domain/application layer from the persistence adapter.
type MatchPort interface {
	CreateMatch(ctx context.Context, match *model.Match) error
	GetMatchByID(ctx context.Context, id uint64) (*model.Match, error)
	GetPaginatedMatches(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Match, int64, error)
	GetMatchesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]model.Match, int64, error)
	GetMatchesByTeamID(ctx context.Context, teamID uint64, sort string, order string, page int, pageSize int) ([]model.Match, int64, error)
	GetNextMatchByTeamID(ctx context.Context, teamID uint64) (*model.Match, error)
	GetDetailedMatchByID(ctx context.Context, id uint64) (*model.Match, error)
	UpdateMatch(ctx context.Context, id uint64, match *model.Match) error
	DeleteMatch(ctx context.Context, id uint64) error
}

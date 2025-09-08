package domain

import "context"

// LineupRepository defines the interface for lineup persistence operations.
// This interface belongs in the domain layer
type LineupRepository interface {
	CreateLineup(ctx context.Context, lineup *Lineup) error
	GetLineupByID(ctx context.Context, id uint64) (*Lineup, error)
	UpdateLineup(ctx context.Context, id uint64, lineup *Lineup) error
	DeleteLineup(ctx context.Context, id uint64) error
	GetPaginatedLineups(ctx context.Context, sort string, order string, page int, pageSize int) ([]Lineup, int64, error)

	// Match-specific operations
	GetLineupsByMatchID(ctx context.Context, matchID uint64) ([]Lineup, error)
	GetStartingLineupsByMatchID(ctx context.Context, matchID uint64) ([]Lineup, error)
	GetSubstitutesLineupsByMatchID(ctx context.Context, matchID uint64) ([]Lineup, error)

	// Player-specific operations
	GetLineupsByPlayerID(ctx context.Context, playerID uint64) ([]Lineup, error)
}

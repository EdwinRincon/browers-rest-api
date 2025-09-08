package domain

import (
	"context"
)

// SeasonRepository defines the interface for season persistence operations.
// This port belongs in the domain layer following hexagonal architecture.
// All methods preserved from existing api/repository and internal/infrastructure interfaces.
type SeasonRepository interface {
	// Standard CRUD methods
	CreateSeason(ctx context.Context, season *Season) error
	GetSeasonByID(ctx context.Context, id uint64) (*Season, error)
	GetSeasonByYear(ctx context.Context, year uint16) (*Season, error)
	GetCurrentSeason(ctx context.Context) (*Season, error)
	GetPaginatedSeasons(ctx context.Context, sort string, order string, page int, pageSize int) ([]Season, int64, error)
	UpdateSeason(ctx context.Context, id uint64, season *Season) error
	DeleteSeason(ctx context.Context, id uint64) error

	// Business logic specific methods
	SetCurrentSeason(ctx context.Context, id uint64) error
}

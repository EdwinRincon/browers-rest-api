package ports

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
)

// SeasonPort defines the interface for season data access operations.
// This port separates the domain/application layer from the persistence adapter.
type SeasonPort interface {
	CreateSeason(ctx context.Context, season *model.Season) error
	GetSeasonByID(ctx context.Context, id uint64) (*model.Season, error)
	GetSeasonByYear(ctx context.Context, year uint16) (*model.Season, error)
	GetCurrentSeason(ctx context.Context) (*model.Season, error)
	GetPaginatedSeasons(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Season, int64, error)
	UpdateSeason(ctx context.Context, id uint64, season *model.Season) error
	DeleteSeason(ctx context.Context, id uint64) error
	SetCurrentSeason(ctx context.Context, id uint64) error
}

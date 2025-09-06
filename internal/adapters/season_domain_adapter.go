package adapters

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

// SeasonDomainAdapter adapts the persistence SeasonPort to work with domain entities.
// This adapter translates between domain entities and persistence models.
type SeasonDomainAdapter struct {
	persistencePort ports.SeasonPort
}

// NewSeasonDomainAdapter creates a new SeasonDomainAdapter.
func NewSeasonDomainAdapter(persistencePort ports.SeasonPort) ports.SeasonDomainPort {
	return &SeasonDomainAdapter{
		persistencePort: persistencePort,
	}
}

func (a *SeasonDomainAdapter) CreateSeason(ctx context.Context, season *domain.Season) error {
	persistenceSeason := mapper.SeasonDomainToPersistence(season)
	return a.persistencePort.CreateSeason(ctx, persistenceSeason)
}

func (a *SeasonDomainAdapter) GetSeasonByID(ctx context.Context, id uint64) (*domain.Season, error) {
	persistenceSeason, err := a.persistencePort.GetSeasonByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if persistenceSeason == nil {
		return nil, nil
	}
	return mapper.SeasonPersistenceToDomain(persistenceSeason), nil
}

func (a *SeasonDomainAdapter) GetSeasonByYear(ctx context.Context, year uint16) (*domain.Season, error) {
	persistenceSeason, err := a.persistencePort.GetSeasonByYear(ctx, year)
	if err != nil {
		return nil, err
	}
	if persistenceSeason == nil {
		return nil, nil
	}
	return mapper.SeasonPersistenceToDomain(persistenceSeason), nil
}

func (a *SeasonDomainAdapter) GetCurrentSeason(ctx context.Context) (*domain.Season, error) {
	persistenceSeason, err := a.persistencePort.GetCurrentSeason(ctx)
	if err != nil {
		return nil, err
	}
	if persistenceSeason == nil {
		return nil, nil
	}
	return mapper.SeasonPersistenceToDomain(persistenceSeason), nil
}

func (a *SeasonDomainAdapter) GetPaginatedSeasons(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Season, int64, error) {
	persistenceSeasons, total, err := a.persistencePort.GetPaginatedSeasons(ctx, sort, order, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	domainSeasons := mapper.SeasonPersistenceListToDomain(persistenceSeasons)
	return domainSeasons, total, nil
}

func (a *SeasonDomainAdapter) UpdateSeason(ctx context.Context, id uint64, season *domain.Season) error {
	persistenceSeason := mapper.SeasonDomainToPersistence(season)
	return a.persistencePort.UpdateSeason(ctx, id, persistenceSeason)
}

func (a *SeasonDomainAdapter) DeleteSeason(ctx context.Context, id uint64) error {
	return a.persistencePort.DeleteSeason(ctx, id)
}

func (a *SeasonDomainAdapter) SetCurrentSeason(ctx context.Context, id uint64) error {
	return a.persistencePort.SetCurrentSeason(ctx, id)
}

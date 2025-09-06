package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// SeasonDomainToPersistence converts a domain Season to a persistence model Season.
func SeasonDomainToPersistence(domainSeason *domain.Season) *model.Season {
	if domainSeason == nil {
		return nil
	}

	return &model.Season{
		ID:        domainSeason.ID,
		Year:      domainSeason.Year,
		StartDate: domainSeason.StartDate,
		EndDate:   domainSeason.EndDate,
		IsCurrent: domainSeason.IsCurrent,
		CreatedAt: domainSeason.CreatedAt,
		UpdatedAt: domainSeason.UpdatedAt,
		// Note: Relations (Matches, Articles, TeamStats, PlayerTeams, PlayerStats)
		// are not mapped as they should be handled separately when needed
	}
}

// SeasonPersistenceToDomain converts a persistence model Season to a domain Season.
func SeasonPersistenceToDomain(persistenceSeason *model.Season) *domain.Season {
	if persistenceSeason == nil {
		return nil
	}

	return &domain.Season{
		ID:        persistenceSeason.ID,
		Year:      persistenceSeason.Year,
		StartDate: persistenceSeason.StartDate,
		EndDate:   persistenceSeason.EndDate,
		IsCurrent: persistenceSeason.IsCurrent,
		CreatedAt: persistenceSeason.CreatedAt,
		UpdatedAt: persistenceSeason.UpdatedAt,
	}
}

// SeasonDomainListToPersistence converts a slice of domain Seasons to persistence model Seasons.
func SeasonDomainListToPersistence(domainSeasons []domain.Season) []model.Season {
	persistenceSeasons := make([]model.Season, len(domainSeasons))
	for i, domainSeason := range domainSeasons {
		persistenceSeasons[i] = *SeasonDomainToPersistence(&domainSeason)
	}
	return persistenceSeasons
}

// SeasonPersistenceListToDomain converts a slice of persistence model Seasons to domain Seasons.
func SeasonPersistenceListToDomain(persistenceSeasons []model.Season) []domain.Season {
	domainSeasons := make([]domain.Season, len(persistenceSeasons))
	for i, persistenceSeason := range persistenceSeasons {
		domainSeasons[i] = *SeasonPersistenceToDomain(&persistenceSeason)
	}
	return domainSeasons
}

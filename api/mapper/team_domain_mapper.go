package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// TeamDomainToPersistence converts a domain Team to a persistence model Team.
func TeamDomainToPersistence(domainTeam *domain.Team) *model.Team {
	if domainTeam == nil {
		return nil
	}

	return &model.Team{
		ID:          domainTeam.ID,
		FullName:    domainTeam.FullName,
		ShortName:   domainTeam.ShortName,
		Color:       domainTeam.Color,
		Color2:      domainTeam.Color2,
		Shield:      domainTeam.Shield,
		NextMatchID: domainTeam.NextMatchID,
		CreatedAt:   domainTeam.CreatedAt,
		UpdatedAt:   domainTeam.UpdatedAt,
		// Note: Relations (NextMatch, PlayerTeams, HomeMatches, AwayMatches, TeamStats, PlayerStats)
		// are not mapped as they should be handled separately when needed
	}
}

// TeamPersistenceToDomain converts a persistence model Team to a domain Team.
func TeamPersistenceToDomain(persistenceTeam *model.Team) *domain.Team {
	if persistenceTeam == nil {
		return nil
	}

	return &domain.Team{
		ID:          persistenceTeam.ID,
		FullName:    persistenceTeam.FullName,
		ShortName:   persistenceTeam.ShortName,
		Color:       persistenceTeam.Color,
		Color2:      persistenceTeam.Color2,
		Shield:      persistenceTeam.Shield,
		NextMatchID: persistenceTeam.NextMatchID,
		CreatedAt:   persistenceTeam.CreatedAt,
		UpdatedAt:   persistenceTeam.UpdatedAt,
	}
}

// TeamDomainListToPersistence converts a slice of domain Teams to persistence model Teams.
func TeamDomainListToPersistence(domainTeams []domain.Team) []model.Team {
	persistenceTeams := make([]model.Team, len(domainTeams))
	for i, domainTeam := range domainTeams {
		persistenceTeams[i] = *TeamDomainToPersistence(&domainTeam)
	}
	return persistenceTeams
}

// TeamPersistenceListToDomain converts a slice of persistence model Teams to domain Teams.
func TeamPersistenceListToDomain(persistenceTeams []model.Team) []domain.Team {
	domainTeams := make([]domain.Team, len(persistenceTeams))
	for i, persistenceTeam := range persistenceTeams {
		domainTeams[i] = *TeamPersistenceToDomain(&persistenceTeam)
	}
	return domainTeams
}

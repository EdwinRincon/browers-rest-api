package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// PlayerDomainToPersistence converts a domain Player to a persistence model Player.
func PlayerDomainToPersistence(domainPlayer *domain.Player) *model.Player {
	if domainPlayer == nil {
		return nil
	}

	return &model.Player{
		ID:            domainPlayer.ID,
		NickName:      domainPlayer.NickName,
		Height:        domainPlayer.Height,
		Country:       domainPlayer.Country,
		Country2:      domainPlayer.Country2,
		Foot:          domainPlayer.Foot,
		Age:           domainPlayer.Age,
		SquadNumber:   domainPlayer.SquadNumber,
		Rating:        domainPlayer.Rating,
		Matches:       domainPlayer.Matches,
		YCards:        domainPlayer.YCards,
		RCards:        domainPlayer.RCards,
		Goals:         domainPlayer.Goals,
		Assists:       domainPlayer.Assists,
		Saves:         domainPlayer.Saves,
		Position:      domainPlayer.Position,
		Injured:       domainPlayer.Injured,
		CareerSummary: domainPlayer.CareerSummary,
		MVPCount:      domainPlayer.MVPCount,
		UserID:        domainPlayer.UserID,
		CreatedAt:     domainPlayer.CreatedAt,
		UpdatedAt:     domainPlayer.UpdatedAt,
		// Note: Relations (PlayerTeams, Lineups, PlayerStats, User) are not mapped
		// as they should be handled separately when needed
	}
}

// PlayerPersistenceToDomain converts a persistence model Player to a domain Player.
func PlayerPersistenceToDomain(persistencePlayer *model.Player) *domain.Player {
	if persistencePlayer == nil {
		return nil
	}

	return &domain.Player{
		ID:            persistencePlayer.ID,
		NickName:      persistencePlayer.NickName,
		Height:        persistencePlayer.Height,
		Country:       persistencePlayer.Country,
		Country2:      persistencePlayer.Country2,
		Foot:          persistencePlayer.Foot,
		Age:           persistencePlayer.Age,
		SquadNumber:   persistencePlayer.SquadNumber,
		Rating:        persistencePlayer.Rating,
		Matches:       persistencePlayer.Matches,
		YCards:        persistencePlayer.YCards,
		RCards:        persistencePlayer.RCards,
		Goals:         persistencePlayer.Goals,
		Assists:       persistencePlayer.Assists,
		Saves:         persistencePlayer.Saves,
		Position:      persistencePlayer.Position,
		Injured:       persistencePlayer.Injured,
		CareerSummary: persistencePlayer.CareerSummary,
		MVPCount:      persistencePlayer.MVPCount,
		UserID:        persistencePlayer.UserID,
		CreatedAt:     persistencePlayer.CreatedAt,
		UpdatedAt:     persistencePlayer.UpdatedAt,
	}
}

// PlayerDomainListToPersistence converts a slice of domain Players to persistence model Players.
func PlayerDomainListToPersistence(domainPlayers []domain.Player) []model.Player {
	persistencePlayers := make([]model.Player, len(domainPlayers))
	for i, domainPlayer := range domainPlayers {
		persistencePlayers[i] = *PlayerDomainToPersistence(&domainPlayer)
	}
	return persistencePlayers
}

// PlayerPersistenceListToDomain converts a slice of persistence model Players to domain Players.
func PlayerPersistenceListToDomain(persistencePlayers []model.Player) []domain.Player {
	domainPlayers := make([]domain.Player, len(persistencePlayers))
	for i, persistencePlayer := range persistencePlayers {
		domainPlayers[i] = *PlayerPersistenceToDomain(&persistencePlayer)
	}
	return domainPlayers
}

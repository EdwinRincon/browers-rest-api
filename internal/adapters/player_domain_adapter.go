package adapters

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

// PlayerDomainAdapter adapts the persistence PlayerRepository to work with domain entities.
// This adapter translates between domain entities and persistence models.
type PlayerDomainAdapter struct {
	persistenceRepo persistence.PlayerRepository
}

// NewPlayerDomainAdapter creates a new PlayerDomainAdapter.
func NewPlayerDomainAdapter(persistenceRepo persistence.PlayerRepository) ports.PlayerPort {
	return &PlayerDomainAdapter{
		persistenceRepo: persistenceRepo,
	}
}

func (a *PlayerDomainAdapter) CreatePlayer(ctx context.Context, player *domain.Player) error {
	persistencePlayer := mapper.PlayerDomainToPersistence(player)
	return a.persistenceRepo.CreatePlayer(ctx, persistencePlayer)
}

func (a *PlayerDomainAdapter) GetPlayerByID(ctx context.Context, id uint64) (*domain.Player, error) {
	persistencePlayer, err := a.persistenceRepo.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if persistencePlayer == nil {
		return nil, nil
	}
	return mapper.PlayerPersistenceToDomain(persistencePlayer), nil
}

func (a *PlayerDomainAdapter) GetPlayerByNickName(ctx context.Context, nickName string) (*domain.Player, error) {
	persistencePlayer, err := a.persistenceRepo.GetPlayerByNickName(ctx, nickName)
	if err != nil {
		return nil, err
	}
	if persistencePlayer == nil {
		return nil, nil
	}
	return mapper.PlayerPersistenceToDomain(persistencePlayer), nil
}

func (a *PlayerDomainAdapter) GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Player, int64, error) {
	persistencePlayers, total, err := a.persistenceRepo.GetPaginatedPlayers(ctx, sort, order, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	domainPlayers := mapper.PlayerPersistenceListToDomain(persistencePlayers)
	return domainPlayers, total, nil
}

func (a *PlayerDomainAdapter) UpdatePlayer(ctx context.Context, id uint64, player *domain.Player) error {
	persistencePlayer := mapper.PlayerDomainToPersistence(player)
	return a.persistenceRepo.UpdatePlayer(ctx, id, persistencePlayer)
}

func (a *PlayerDomainAdapter) DeletePlayer(ctx context.Context, id uint64) error {
	return a.persistenceRepo.DeletePlayer(ctx, id)
}

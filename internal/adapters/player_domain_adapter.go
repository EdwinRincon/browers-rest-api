package adapters

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

// PlayerDomainAdapter adapts the persistence PlayerPort to work with domain entities.
// This adapter translates between domain entities and persistence models.
type PlayerDomainAdapter struct {
	persistencePort ports.PlayerPort
}

// NewPlayerDomainAdapter creates a new PlayerDomainAdapter.
func NewPlayerDomainAdapter(persistencePort ports.PlayerPort) ports.PlayerDomainPort {
	return &PlayerDomainAdapter{
		persistencePort: persistencePort,
	}
}

func (a *PlayerDomainAdapter) CreatePlayer(ctx context.Context, player *domain.Player) error {
	persistencePlayer := mapper.PlayerDomainToPersistence(player)
	return a.persistencePort.CreatePlayer(ctx, persistencePlayer)
}

func (a *PlayerDomainAdapter) GetPlayerByID(ctx context.Context, id uint64) (*domain.Player, error) {
	persistencePlayer, err := a.persistencePort.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if persistencePlayer == nil {
		return nil, nil
	}
	return mapper.PlayerPersistenceToDomain(persistencePlayer), nil
}

func (a *PlayerDomainAdapter) GetPlayerByNickName(ctx context.Context, nickName string) (*domain.Player, error) {
	persistencePlayer, err := a.persistencePort.GetPlayerByNickName(ctx, nickName)
	if err != nil {
		return nil, err
	}
	if persistencePlayer == nil {
		return nil, nil
	}
	return mapper.PlayerPersistenceToDomain(persistencePlayer), nil
}

func (a *PlayerDomainAdapter) GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Player, int64, error) {
	persistencePlayers, total, err := a.persistencePort.GetPaginatedPlayers(ctx, sort, order, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	domainPlayers := mapper.PlayerPersistenceListToDomain(persistencePlayers)
	return domainPlayers, total, nil
}

func (a *PlayerDomainAdapter) UpdatePlayer(ctx context.Context, id uint64, player *domain.Player) error {
	persistencePlayer := mapper.PlayerDomainToPersistence(player)
	return a.persistencePort.UpdatePlayer(ctx, id, persistencePlayer)
}

func (a *PlayerDomainAdapter) DeletePlayer(ctx context.Context, id uint64) error {
	return a.persistencePort.DeletePlayer(ctx, id)
}

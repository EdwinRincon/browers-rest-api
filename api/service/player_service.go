package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type PlayerService interface {
	CreatePlayer(ctx context.Context, player *model.Players) error
	GetPlayerByID(ctx context.Context, id uint64) (*model.Players, error)
	ListPlayers(ctx context.Context, page uint64) ([]*model.Players, error)
	UpdatePlayer(ctx context.Context, player *model.Players) error
	DeletePlayer(ctx context.Context, id uint64) error
}

type playerService struct {
	PlayerRepository repository.PlayerRepository
}

func NewPlayerService(playerRepo repository.PlayerRepository) PlayerService {
	return &playerService{
		PlayerRepository: playerRepo,
	}
}

func (s *playerService) CreatePlayer(ctx context.Context, player *model.Players) error {
	return s.PlayerRepository.CreatePlayer(ctx, player)
}

func (s *playerService) GetPlayerByID(ctx context.Context, id uint64) (*model.Players, error) {
	return s.PlayerRepository.GetPlayerByID(ctx, id)
}

func (s *playerService) ListPlayers(ctx context.Context, page uint64) ([]*model.Players, error) {
	return s.PlayerRepository.ListPlayers(ctx, page)
}

func (s *playerService) UpdatePlayer(ctx context.Context, player *model.Players) error {
	return s.PlayerRepository.UpdatePlayer(ctx, player)
}

func (s *playerService) DeletePlayer(ctx context.Context, id uint64) error {
	return s.PlayerRepository.DeletePlayer(ctx, id)
}

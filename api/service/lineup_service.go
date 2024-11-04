package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type LineupService interface {
	CreateLineup(ctx context.Context, lineup *model.Lineups) error
	GetLineupByID(ctx context.Context, id uint64) (*model.Lineups, error)
	ListLineups(ctx context.Context, page uint64) ([]*model.Lineups, error)
	UpdateLineup(ctx context.Context, lineup *model.Lineups) error
	DeleteLineup(ctx context.Context, id uint64) error
	GetLineupsByMatch(ctx context.Context, matchID uint64) ([]*model.Lineups, error)
}

type lineupService struct {
	LineupRepository repository.LineupRepository
}

func NewLineupService(lineupRepo repository.LineupRepository) LineupService {
	return &lineupService{
		LineupRepository: lineupRepo,
	}
}

func (s *lineupService) CreateLineup(ctx context.Context, lineup *model.Lineups) error {
	return s.LineupRepository.CreateLineup(ctx, lineup)
}

func (s *lineupService) GetLineupByID(ctx context.Context, id uint64) (*model.Lineups, error) {
	return s.LineupRepository.GetLineupByID(ctx, id)
}

func (s *lineupService) ListLineups(ctx context.Context, page uint64) ([]*model.Lineups, error) {
	return s.LineupRepository.ListLineups(ctx, page)
}

func (s *lineupService) UpdateLineup(ctx context.Context, lineup *model.Lineups) error {
	return s.LineupRepository.UpdateLineup(ctx, lineup)
}

func (s *lineupService) DeleteLineup(ctx context.Context, id uint64) error {
	return s.LineupRepository.DeleteLineup(ctx, id)
}

func (s *lineupService) GetLineupsByMatch(ctx context.Context, matchID uint64) ([]*model.Lineups, error) {
	return s.LineupRepository.GetLineupsByMatch(ctx, matchID)
}

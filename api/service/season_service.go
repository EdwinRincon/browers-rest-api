package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type SeasonService interface {
	CreateSeason(ctx context.Context, season *model.Season) error
	GetSeasonByID(ctx context.Context, id uint8) (*model.Season, error)
	GetAllSeasons(ctx context.Context, page uint64) ([]*model.Season, error)
	UpdateSeason(ctx context.Context, season *model.Season) error
	DeleteSeason(ctx context.Context, id uint8) error
}

type seasonService struct {
	SeasonRepository repository.SeasonRepository
}

func NewSeasonService(seasonRepo repository.SeasonRepository) SeasonService {
	return &seasonService{
		SeasonRepository: seasonRepo,
	}
}

func (s *seasonService) CreateSeason(ctx context.Context, season *model.Season) error {
	return s.SeasonRepository.CreateSeason(ctx, season)
}

func (s *seasonService) GetSeasonByID(ctx context.Context, id uint8) (*model.Season, error) {
	return s.SeasonRepository.GetSeasonByID(ctx, id)
}

func (s *seasonService) GetAllSeasons(ctx context.Context, page uint64) ([]*model.Season, error) {
	return s.SeasonRepository.GetAllSeasons(ctx, page)
}

func (s *seasonService) UpdateSeason(ctx context.Context, season *model.Season) error {
	return s.SeasonRepository.UpdateSeason(ctx, season)
}

func (s *seasonService) DeleteSeason(ctx context.Context, id uint8) error {
	return s.SeasonRepository.DeleteSeason(ctx, id)
}

package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

const errFailedToGetSeason = "failed to get season: %w"

type SeasonService interface {
	CreateSeason(ctx context.Context, createRequest *dto.CreateSeasonRequest) (*dto.SeasonResponse, error)
	GetSeasonByID(ctx context.Context, id uint64) (*model.Season, error)
	GetSeasonByYear(ctx context.Context, year uint16) (*model.Season, error)
	GetCurrentSeason(ctx context.Context) (*model.Season, error)
	GetPaginatedSeasons(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Season, int64, error)
	UpdateSeason(ctx context.Context, id uint64, updateRequest *dto.UpdateSeasonRequest) (*dto.SeasonResponse, error)
	DeleteSeason(ctx context.Context, id uint64) error
	SetCurrentSeason(ctx context.Context, id uint64) error
}

type seasonService struct {
	SeasonRepository repository.SeasonRepository
}

func NewSeasonService(seasonRepo repository.SeasonRepository) SeasonService {
	return &seasonService{
		SeasonRepository: seasonRepo,
	}
}

func (s *seasonService) CreateSeason(ctx context.Context, createRequest *dto.CreateSeasonRequest) (*dto.SeasonResponse, error) {
	// Check if season with same year already exists
	existingSeason, err := s.SeasonRepository.GetSeasonByYear(ctx, createRequest.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing season: %w", err)
	}
	if existingSeason != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	// Convert DTO to model
	season := mapper.ToSeason(createRequest)

	if err := s.SeasonRepository.CreateSeason(ctx, season); err != nil {
		return nil, fmt.Errorf("failed to create season: %w", err)
	}

	return mapper.ToSeasonResponse(season), nil
}

func (s *seasonService) GetSeasonByID(ctx context.Context, id uint64) (*model.Season, error) {
	season, err := s.SeasonRepository.GetActiveSeasonByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get season by ID: %w", err)
	}
	if season == nil {
		return nil, constants.ErrRecordNotFound
	}
	return season, nil
}

func (s *seasonService) GetSeasonByYear(ctx context.Context, year uint16) (*model.Season, error) {
	season, err := s.SeasonRepository.GetSeasonByYear(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get season by year: %w", err)
	}
	if season == nil {
		return nil, constants.ErrRecordNotFound
	}
	return season, nil
}

func (s *seasonService) GetCurrentSeason(ctx context.Context) (*model.Season, error) {
	season, err := s.SeasonRepository.GetCurrentSeason(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current season: %w", err)
	}
	if season == nil {
		return nil, constants.ErrRecordNotFound
	}
	return season, nil
}

func (s *seasonService) GetPaginatedSeasons(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Season, int64, error) {
	return s.SeasonRepository.GetPaginatedSeasons(ctx, sort, order, page, pageSize)
}

func (s *seasonService) UpdateSeason(ctx context.Context, id uint64, updateRequest *dto.UpdateSeasonRequest) (*dto.SeasonResponse, error) {
	// Fetch existing season
	season, err := s.SeasonRepository.GetActiveSeasonByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(errFailedToGetSeason, err)
	}
	if season == nil {
		return nil, constants.ErrRecordNotFound
	}

	// If updating the year, check for uniqueness
	if updateRequest.Year != nil && *updateRequest.Year != season.Year {
		existingSeason, err := s.SeasonRepository.GetSeasonByYear(ctx, *updateRequest.Year)
		if err != nil {
			return nil, fmt.Errorf("failed to check year uniqueness: %w", err)
		}
		if existingSeason != nil && existingSeason.ID != id {
			return nil, constants.ErrRecordAlreadyExists
		}
	}

	// Update fields from DTO
	mapper.UpdateSeasonFromDTO(season, updateRequest)

	if err := s.SeasonRepository.UpdateSeason(ctx, id, season); err != nil {
		return nil, fmt.Errorf("failed to update season: %w", err)
	}

	return mapper.ToSeasonResponse(season), nil
}

func (s *seasonService) DeleteSeason(ctx context.Context, id uint64) error {
	// Check if season exists
	season, err := s.SeasonRepository.GetActiveSeasonByID(ctx, id)
	if err != nil {
		return fmt.Errorf(errFailedToGetSeason, err)
	}
	if season == nil {
		return constants.ErrRecordNotFound
	}

	// Do not allow deletion of current season
	if season.IsCurrent {
		return fmt.Errorf("cannot delete current season")
	}

	return s.SeasonRepository.DeleteSeason(ctx, id)
}

func (s *seasonService) SetCurrentSeason(ctx context.Context, id uint64) error {
	// Check if season exists
	season, err := s.SeasonRepository.GetActiveSeasonByID(ctx, id)
	if err != nil {
		return fmt.Errorf(errFailedToGetSeason, err)
	}
	if season == nil {
		return constants.ErrRecordNotFound
	}

	return s.SeasonRepository.SetCurrentSeason(ctx, id)
}

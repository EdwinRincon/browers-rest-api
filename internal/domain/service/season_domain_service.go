package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

// SeasonDomainService implements business logic for Season operations.
// It contains domain rules and validation while being infrastructure-agnostic.
type SeasonDomainService struct {
	seasonPort ports.SeasonDomainPort
}

// NewSeasonDomainService creates a new SeasonDomainService instance.
func NewSeasonDomainService(seasonPort ports.SeasonDomainPort) *SeasonDomainService {
	return &SeasonDomainService{
		seasonPort: seasonPort,
	}
}

// CreateSeason creates a new season with business rule validation.
func (s *SeasonDomainService) CreateSeason(ctx context.Context, season *domain.Season) error {
	// Validate domain rules
	if !season.IsValid() {
		return fmt.Errorf("season fails domain validation")
	}

	// Business rule: Check if season with same year already exists
	existingSeason, err := s.seasonPort.GetSeasonByYear(ctx, season.Year)
	if err != nil {
		return fmt.Errorf("failed to check existing season: %w", err)
	}
	if existingSeason != nil {
		return constants.ErrRecordAlreadyExists
	}

	// Business rule: If this season is set as current, ensure no overlapping current seasons
	if season.IsCurrent {
		currentSeason, err := s.seasonPort.GetCurrentSeason(ctx)
		if err != nil {
			return fmt.Errorf("failed to check current season: %w", err)
		}
		if currentSeason != nil {
			return fmt.Errorf("season %d is already set as current", currentSeason.Year)
		}
	}

	return s.seasonPort.CreateSeason(ctx, season)
}

// GetSeasonByID retrieves a season by ID.
func (s *SeasonDomainService) GetSeasonByID(ctx context.Context, id uint64) (*domain.Season, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid season ID")
	}

	season, err := s.seasonPort.GetSeasonByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get season by ID: %w", err)
	}

	if season == nil {
		return nil, constants.ErrRecordNotFound
	}

	return season, nil
}

// GetSeasonByYear retrieves a season by year.
func (s *SeasonDomainService) GetSeasonByYear(ctx context.Context, year uint16) (*domain.Season, error) {
	if year == 0 {
		return nil, fmt.Errorf("invalid season year")
	}

	season, err := s.seasonPort.GetSeasonByYear(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get season by year: %w", err)
	}

	if season == nil {
		return nil, fmt.Errorf("season for year %d not found", year)
	}

	return season, nil
}

// GetCurrentSeason retrieves the current active season.
func (s *SeasonDomainService) GetCurrentSeason(ctx context.Context) (*domain.Season, error) {
	season, err := s.seasonPort.GetCurrentSeason(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current season: %w", err)
	}

	if season == nil {
		return nil, constants.ErrRecordNotFound
	}

	return season, nil
}

// GetPaginatedSeasons retrieves paginated seasons.
func (s *SeasonDomainService) GetPaginatedSeasons(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Season, int64, error) {
	// Validate pagination parameters
	if page < 0 {
		return nil, 0, fmt.Errorf("page must be non-negative")
	}
	if pageSize <= 0 || pageSize > 100 {
		return nil, 0, fmt.Errorf("page size must be between 1 and 100")
	}

	return s.seasonPort.GetPaginatedSeasons(ctx, sort, order, page, pageSize)
}

// UpdateSeason updates an existing season with business rule validation.
func (s *SeasonDomainService) UpdateSeason(ctx context.Context, id uint64, season *domain.Season) error {
	if id == 0 {
		return fmt.Errorf("invalid season ID")
	}

	// Validate domain rules
	if !season.IsValid() {
		return fmt.Errorf("season fails domain validation")
	}

	// Check if season exists
	existingSeason, err := s.seasonPort.GetSeasonByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get existing season: %w", err)
	}
	if existingSeason == nil {
		return constants.ErrRecordNotFound
	}

	// Business rule: If updating year, check for uniqueness
	if season.Year != existingSeason.Year {
		conflictingSeason, err := s.seasonPort.GetSeasonByYear(ctx, season.Year)
		if err != nil {
			return fmt.Errorf("failed to check year uniqueness: %w", err)
		}
		if conflictingSeason != nil && conflictingSeason.ID != id {
			return constants.ErrRecordAlreadyExists
		}
	}

	// Business rule: If setting as current, ensure no other season is current
	if season.IsCurrent && !existingSeason.IsCurrent {
		currentSeason, err := s.seasonPort.GetCurrentSeason(ctx)
		if err != nil {
			return fmt.Errorf("failed to check current season: %w", err)
		}
		if currentSeason != nil && currentSeason.ID != id {
			return fmt.Errorf("season %d is already set as current", currentSeason.Year)
		}
	}

	return s.seasonPort.UpdateSeason(ctx, id, season)
}

// DeleteSeason deletes a season with business rule validation.
func (s *SeasonDomainService) DeleteSeason(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid season ID")
	}

	// Check if season exists
	season, err := s.seasonPort.GetSeasonByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get season: %w", err)
	}
	if season == nil {
		return fmt.Errorf("season with ID %d not found", id)
	}

	// Business rule: Cannot delete current season
	if season.IsCurrent {
		return fmt.Errorf("cannot delete current season")
	}

	return s.seasonPort.DeleteSeason(ctx, id)
}

// SetCurrentSeason sets a season as the current one.
func (s *SeasonDomainService) SetCurrentSeason(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid season ID")
	}

	// Check if season exists
	season, err := s.seasonPort.GetSeasonByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get season: %w", err)
	}
	if season == nil {
		return fmt.Errorf("season with ID %d not found", id)
	}

	// Business rule: Season must be valid to be set as current
	if !season.IsValid() {
		return fmt.Errorf("cannot set invalid season as current")
	}

	return s.seasonPort.SetCurrentSeason(ctx, id)
}

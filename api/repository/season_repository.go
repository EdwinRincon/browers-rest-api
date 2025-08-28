package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

const (
	errClearCurrentSeasonFlag = "failed to clear current season flag: %w"
	whereID                   = "id = ?"
	whereSeasonID             = "season_id = ?"
)

type SeasonRepository interface {
	CreateSeason(ctx context.Context, season *model.Season) error
	GetActiveSeasonByID(ctx context.Context, id uint64) (*model.Season, error)
	GetUnscopedSeasonByID(ctx context.Context, id uint64) (*model.Season, error)
	GetSeasonByYear(ctx context.Context, year uint16) (*model.Season, error)
	GetCurrentSeason(ctx context.Context) (*model.Season, error)
	GetPaginatedSeasons(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Season, int64, error)
	UpdateSeason(ctx context.Context, id uint64, season *model.Season) error
	DeleteSeason(ctx context.Context, id uint64) error
	SetCurrentSeason(ctx context.Context, id uint64) error
}

type SeasonRepositoryImpl struct {
	db *gorm.DB
}

func NewSeasonRepository(db *gorm.DB) SeasonRepository {
	return &SeasonRepositoryImpl{db: db}
}

func (sr *SeasonRepositoryImpl) CreateSeason(ctx context.Context, season *model.Season) error {
	// If this season is set as current, make sure no other season is current
	if season.IsCurrent {
		if err := sr.clearCurrentSeasons(ctx); err != nil {
			return fmt.Errorf(errClearCurrentSeasonFlag, err)
		}
	}

	return sr.db.WithContext(ctx).Create(season).Error
}

// clearCurrentSeasons sets IsCurrent=false for all seasons
func (sr *SeasonRepositoryImpl) clearCurrentSeasons(ctx context.Context) error {
	return sr.db.WithContext(ctx).Model(&model.Season{}).Where("is_current = ?", true).Update("is_current", false).Error
}

func (sr *SeasonRepositoryImpl) GetActiveSeasonByID(ctx context.Context, id uint64) (*model.Season, error) {
	var season model.Season
	result := sr.db.WithContext(ctx).
		Preload("Matches").
		Preload("Articles").
		Preload("TeamStats").
		Preload("PlayerTeams").
		Preload("PlayerStats").
		Where(whereID, id).
		First(&season)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error getting season by ID: %w", result.Error)
	}

	return &season, nil
}

func (sr *SeasonRepositoryImpl) GetUnscopedSeasonByID(ctx context.Context, id uint64) (*model.Season, error) {
	var season model.Season
	result := sr.db.WithContext(ctx).
		Unscoped().
		Preload("Matches").
		Preload("Articles").
		Preload("TeamStats").
		Preload("PlayerTeams").
		Preload("PlayerStats").
		Where(whereID, id).
		First(&season)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error getting season by ID: %w", result.Error)
	}

	return &season, nil
}

func (sr *SeasonRepositoryImpl) GetSeasonByYear(ctx context.Context, year uint16) (*model.Season, error) {
	var season model.Season
	result := sr.db.WithContext(ctx).
		Where("year = ?", year).
		First(&season)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error getting season by year: %w", result.Error)
	}

	return &season, nil
}

func (sr *SeasonRepositoryImpl) GetCurrentSeason(ctx context.Context) (*model.Season, error) {
	var season model.Season
	result := sr.db.WithContext(ctx).
		Where("is_current = ?", true).
		First(&season)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error getting current season: %w", result.Error)
	}

	return &season, nil
}

func (sr *SeasonRepositoryImpl) GetPaginatedSeasons(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Season, int64, error) {
	var seasons []model.Season
	var total int64

	// Count total records
	countQuery := sr.db.WithContext(ctx).Model(&model.Season{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total seasons: %w", err)
	}

	// Build the data query
	query := sr.db.WithContext(ctx).Model(&model.Season{})

	// Apply sorting if provided
	if sort != "" && (order == "asc" || order == "desc") {
		query = query.Order(fmt.Sprintf("%s %s", sort, order))
	} else {
		// Default ordering by year in descending order
		query = query.Order("year desc")
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&seasons).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching seasons: %w", err)
	}

	return seasons, total, nil
}

func (sr *SeasonRepositoryImpl) UpdateSeason(ctx context.Context, id uint64, season *model.Season) error {
	// If this season is being set as current, clear current flag from other seasons
	if season.IsCurrent {
		if err := sr.clearCurrentSeasons(ctx); err != nil {
			return fmt.Errorf(errClearCurrentSeasonFlag, err)
		}
	}

	return sr.db.WithContext(ctx).
		Model(&model.Season{}).
		Where(whereID, id).
		Select("*").
		Updates(season).Error
}

func (sr *SeasonRepositoryImpl) DeleteSeason(ctx context.Context, id uint64) error {
	return sr.db.WithContext(ctx).Delete(&model.Season{}, id).Error
}

func (sr *SeasonRepositoryImpl) SetCurrentSeason(ctx context.Context, id uint64) error {
	// First, clear current flag from all seasons
	if err := sr.clearCurrentSeasons(ctx); err != nil {
		return fmt.Errorf("failed to clear current season flag: %w", err)
	}

	// Then set the specified season as current
	result := sr.db.WithContext(ctx).
		Model(&model.Season{}).
		Where(whereID, id).
		Update("is_current", true)

	if result.Error != nil {
		return fmt.Errorf("failed to set current season: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("season with id %d not found", id)
	}

	return nil
}

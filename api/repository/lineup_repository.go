package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

// LineupRepository defines the interface for lineup data access.
type LineupRepository interface {
	CreateLineup(ctx context.Context, lineup *model.Lineup) error
	GetLineupByID(ctx context.Context, id uint64) (*model.Lineup, error)
	GetLineupsByMatchID(ctx context.Context, matchID uint64) ([]model.Lineup, error)
	GetLineupsByPlayerID(ctx context.Context, playerID uint64) ([]model.Lineup, error)
	GetPaginatedLineups(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Lineup, int64, error)
	UpdateLineup(ctx context.Context, id uint64, lineup *model.Lineup) error
	DeleteLineup(ctx context.Context, id uint64) error
	GetStartingLineupsByMatchID(ctx context.Context, matchID uint64) ([]model.Lineup, error)
	GetSubstitutesLineupsByMatchID(ctx context.Context, matchID uint64) ([]model.Lineup, error)
}

// LineupRepositoryImpl implements the LineupRepository interface using GORM.
type LineupRepositoryImpl struct {
	db *gorm.DB
}

// NewLineupRepository creates a new LineupRepository instance.
func NewLineupRepository(db *gorm.DB) LineupRepository {
	return &LineupRepositoryImpl{db: db}
}

// CreateLineup adds a new lineup to the database
func (lr *LineupRepositoryImpl) CreateLineup(ctx context.Context, lineup *model.Lineup) error {
	return lr.db.WithContext(ctx).Create(lineup).Error
}

// GetLineupByID retrieves a lineup by its ID, preloading related entities
func (lr *LineupRepositoryImpl) GetLineupByID(ctx context.Context, id uint64) (*model.Lineup, error) {
	var lineup model.Lineup
	result := lr.db.WithContext(ctx).
		Preload("Player").
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Where("id = ?", id).
		First(&lineup)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &lineup, result.Error
}

// GetLineupsByMatchID retrieves all lineups for a specific match
func (lr *LineupRepositoryImpl) GetLineupsByMatchID(ctx context.Context, matchID uint64) ([]model.Lineup, error) {
	var lineups []model.Lineup
	result := lr.db.WithContext(ctx).
		Preload("Player").
		Where("match_id = ?", matchID).
		Find(&lineups)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting lineups by match ID: %w", result.Error)
	}

	return lineups, nil
}

// GetLineupsByPlayerID retrieves all lineups for a specific player
func (lr *LineupRepositoryImpl) GetLineupsByPlayerID(ctx context.Context, playerID uint64) ([]model.Lineup, error) {
	var lineups []model.Lineup
	result := lr.db.WithContext(ctx).
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Where("player_id = ?", playerID).
		Find(&lineups)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting lineups by player ID: %w", result.Error)
	}

	return lineups, nil
}

// GetPaginatedLineups retrieves a paginated list of lineups with their related entities
func (lr *LineupRepositoryImpl) GetPaginatedLineups(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Lineup, int64, error) {
	var lineups []model.Lineup
	var total int64

	// Count total records
	countQuery := lr.db.WithContext(ctx).Model(&model.Lineup{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total lineups: %w", err)
	}

	// Build the data query with eager loading
	query := lr.db.WithContext(ctx).Model(&model.Lineup{}).
		Preload("Player").
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam")

	// Apply sorting if provided
	if sort != "" && (order == "asc" || order == "desc") {
		query = query.Order(fmt.Sprintf("%s %s", sort, order))
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&lineups).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching lineups: %w", err)
	}

	return lineups, total, nil
}

// UpdateLineup updates an existing lineup
func (lr *LineupRepositoryImpl) UpdateLineup(ctx context.Context, id uint64, lineup *model.Lineup) error {
	return lr.db.WithContext(ctx).
		Model(&model.Lineup{}).
		Where("id = ?", id).
		Select("*").
		Updates(lineup).Error
}

// DeleteLineup removes a lineup by its ID
func (lr *LineupRepositoryImpl) DeleteLineup(ctx context.Context, id uint64) error {
	return lr.db.WithContext(ctx).Delete(&model.Lineup{}, "id = ?", id).Error
}

// GetStartingLineupsByMatchID retrieves all starting lineups for a specific match
func (lr *LineupRepositoryImpl) GetStartingLineupsByMatchID(ctx context.Context, matchID uint64) ([]model.Lineup, error) {
	var lineups []model.Lineup
	result := lr.db.WithContext(ctx).
		Preload("Player").
		Where("match_id = ? AND starting = ?", matchID, true).
		Find(&lineups)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting starting lineups by match ID: %w", result.Error)
	}

	return lineups, nil
}

// GetSubstitutesLineupsByMatchID retrieves all substitute lineups for a specific match
func (lr *LineupRepositoryImpl) GetSubstitutesLineupsByMatchID(ctx context.Context, matchID uint64) ([]model.Lineup, error) {
	var lineups []model.Lineup
	result := lr.db.WithContext(ctx).
		Preload("Player").
		Where("match_id = ? AND starting = ?", matchID, false).
		Find(&lineups)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting substitute lineups by match ID: %w", result.Error)
	}

	return lineups, nil
}

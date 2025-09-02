package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type PlayerStatsRepository interface {
	CreatePlayerStat(ctx context.Context, playerStat *model.PlayerStat) error
	GetPlayerStatByID(ctx context.Context, id uint64) (*model.PlayerStat, error)
	GetPlayerStatsByPlayerID(ctx context.Context, playerID uint64) ([]model.PlayerStat, error)
	GetPlayerStatsByMatchID(ctx context.Context, matchID uint64) ([]model.PlayerStat, error)
	GetPlayerStatsBySeasonID(ctx context.Context, seasonID uint64) ([]model.PlayerStat, error)
	GetPaginatedPlayerStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.PlayerStat, int64, error)
	UpdatePlayerStat(ctx context.Context, id uint64, playerStat *model.PlayerStat) error
	DeletePlayerStat(ctx context.Context, id uint64) error
}

type PlayerStatsRepositoryImpl struct {
	db *gorm.DB
}

func NewPlayerStatsRepository(db *gorm.DB) PlayerStatsRepository {
	return &PlayerStatsRepositoryImpl{db: db}
}

func (psr *PlayerStatsRepositoryImpl) CreatePlayerStat(ctx context.Context, playerStat *model.PlayerStat) error {
	return psr.db.WithContext(ctx).Create(playerStat).Error
}

func (psr *PlayerStatsRepositoryImpl) GetPlayerStatByID(ctx context.Context, id uint64) (*model.PlayerStat, error) {
	var playerStat model.PlayerStat
	result := psr.db.WithContext(ctx).
		Preload("Player").
		Preload("Match").
		Preload("Season").
		Preload("Team").
		Where("id = ?", id).
		First(&playerStat)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting player stat by ID: %w", result.Error)
	}
	return &playerStat, nil
}

// GetPlayerStatsByPlayerID retrieves player stats for a specific player.
func (psr *PlayerStatsRepositoryImpl) GetPlayerStatsByPlayerID(ctx context.Context, playerID uint64) ([]model.PlayerStat, error) {
	var playerStats []model.PlayerStat
	result := psr.db.WithContext(ctx).
		Preload("Player").
		Preload("Match").
		Preload("Season").
		Preload("Team").
		Where("player_id = ?", playerID).
		Find(&playerStats)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting player stats by player ID: %w", result.Error)
	}
	return playerStats, nil
}

// GetPlayerStatsByMatchID retrieves player stats for a specific match.
func (psr *PlayerStatsRepositoryImpl) GetPlayerStatsByMatchID(ctx context.Context, matchID uint64) ([]model.PlayerStat, error) {
	var playerStats []model.PlayerStat
	result := psr.db.WithContext(ctx).
		Preload("Player").
		Preload("Match").
		Preload("Season").
		Preload("Team").
		Where("match_id = ?", matchID).
		Find(&playerStats)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting player stats by match ID: %w", result.Error)
	}
	return playerStats, nil
}

// GetPlayerStatsBySeasonID retrieves player stats for a specific season.
func (psr *PlayerStatsRepositoryImpl) GetPlayerStatsBySeasonID(ctx context.Context, seasonID uint64) ([]model.PlayerStat, error) {
	var playerStats []model.PlayerStat
	result := psr.db.WithContext(ctx).
		Preload("Player").
		Preload("Match").
		Preload("Team").
		Preload("Season").
		Where("season_id = ?", seasonID).
		Find(&playerStats)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting player stats by season ID: %w", result.Error)
	}
	return playerStats, nil
}

// GetPaginatedPlayerStats retrieves a paginated list of player stats.
func (psr *PlayerStatsRepositoryImpl) GetPaginatedPlayerStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.PlayerStat, int64, error) {
	var playerStats []model.PlayerStat
	var total int64

	// Count total records
	countQuery := psr.db.WithContext(ctx).Model(&model.PlayerStat{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total player stats: %w", err)
	}

	// Build the data query with eager loading
	query := psr.db.WithContext(ctx).Model(&model.PlayerStat{}).
		Preload("Player").
		Preload("Match").
		Preload("Season").
		Preload("Team")

	// Apply sorting if provided
	if sort != "" && (order == "asc" || order == "desc") {
		query = query.Order(fmt.Sprintf("%s %s", sort, order))
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&playerStats).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching player stats: %w", err)
	}

	return playerStats, total, nil
}

func (psr *PlayerStatsRepositoryImpl) UpdatePlayerStat(ctx context.Context, id uint64, playerStat *model.PlayerStat) error {
	return psr.db.WithContext(ctx).
		Model(&model.PlayerStat{}).
		Where("id = ?", id).
		Select("*").
		Updates(playerStat).Error
}

func (psr *PlayerStatsRepositoryImpl) DeletePlayerStat(ctx context.Context, id uint64) error {
	return psr.db.WithContext(ctx).Delete(&model.PlayerStat{}, "id = ?", id).Error
}

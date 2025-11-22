package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/adapter/persistence"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"github.com/EdwinRincon/browersfc-api/domain"
	"gorm.io/gorm"
)

type PlayerStatsRepositoryImpl struct {
	db     *gorm.DB
	mapper *persistence.PlayerStatsPersistenceMapper
}

func NewPlayerStatsRepository(db *gorm.DB) *PlayerStatsRepositoryImpl {
	return &PlayerStatsRepositoryImpl{
		db:     db,
		mapper: persistence.NewPlayerStatsPersistenceMapper(),
	}
}

func (psr *PlayerStatsRepositoryImpl) CreatePlayerStat(ctx context.Context, playerStat *domain.PlayerStat) error {
	modelPlayerStat := psr.mapper.DomainToModel(playerStat)
	return psr.db.WithContext(ctx).Create(modelPlayerStat).Error
}

func (psr *PlayerStatsRepositoryImpl) GetPlayerStatByID(ctx context.Context, id uint64) (*domain.PlayerStat, error) {
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

	return psr.mapper.ModelToDomain(&playerStat), nil
}

// GetPlayerStatsByPlayerID retrieves player stats for a specific player.
func (psr *PlayerStatsRepositoryImpl) GetPlayerStatsByPlayerID(ctx context.Context, playerID uint64) ([]domain.PlayerStat, error) {
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
	return psr.mapper.ModelListToDomain(playerStats), nil
}

// GetPlayerStatsByMatchID retrieves player stats for a specific match.
func (psr *PlayerStatsRepositoryImpl) GetPlayerStatsByMatchID(ctx context.Context, matchID uint64) ([]domain.PlayerStat, error) {
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
	return psr.mapper.ModelListToDomain(playerStats), nil
}

// GetPlayerStatsBySeasonID retrieves player stats for a specific season.
func (psr *PlayerStatsRepositoryImpl) GetPlayerStatsBySeasonID(ctx context.Context, seasonID uint64) ([]domain.PlayerStat, error) {
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
	return psr.mapper.ModelListToDomain(playerStats), nil
}

// GetPaginatedPlayerStats retrieves a paginated list of player stats.
func (psr *PlayerStatsRepositoryImpl) GetPaginatedPlayerStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.PlayerStat, int64, error) {
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

	// Apply sorting (safe and validated)
	col, raw, err := BuildOrderClause(EntityPlayerStats, sort, order)
	if err != nil {
		return nil, 0, fmt.Errorf("error building sort clause: %w", err)
	}

	if raw != "" {
		query = query.Order(raw)
	} else {
		query = query.Order(col)
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&playerStats).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching player stats: %w", err)
	}

	return psr.mapper.ModelListToDomain(playerStats), total, nil
}

func (psr *PlayerStatsRepositoryImpl) UpdatePlayerStat(ctx context.Context, id uint64, playerStat *domain.PlayerStat) error {
	modelPlayerStat := psr.mapper.DomainToModel(playerStat)
	return psr.db.WithContext(ctx).
		Model(&model.PlayerStat{}).
		Where("id = ?", id).
		Select("*").
		Updates(modelPlayerStat).Error
}

func (psr *PlayerStatsRepositoryImpl) DeletePlayerStat(ctx context.Context, id uint64) error {
	return psr.db.WithContext(ctx).Delete(&model.PlayerStat{}, "id = ?", id).Error
}

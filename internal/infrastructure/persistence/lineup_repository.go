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

type LineupRepositoryImpl struct {
	db     *gorm.DB
	mapper *persistence.LineupPersistenceMapper
}

func NewLineupRepository(db *gorm.DB) domain.LineupRepository {
	return &LineupRepositoryImpl{
		db:     db,
		mapper: persistence.NewLineupPersistenceMapper(),
	}
}

func (r *LineupRepositoryImpl) CreateLineup(ctx context.Context, lineup *domain.Lineup) error {
	lineupModel := r.mapper.DomainToModel(lineup)
	return r.db.WithContext(ctx).Create(lineupModel).Error
}

func (r *LineupRepositoryImpl) GetLineupByID(ctx context.Context, id uint64) (*domain.Lineup, error) {
	var lineupModel model.Lineup
	result := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Where("id = ?", id).
		First(&lineupModel)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return r.mapper.ModelToDomain(&lineupModel), nil
}

func (r *LineupRepositoryImpl) GetLineupsByMatchID(ctx context.Context, matchID uint64) ([]domain.Lineup, error) {
	var lineupModels []model.Lineup
	result := r.db.WithContext(ctx).
		Preload("Player").
		Where("match_id = ?", matchID).
		Find(&lineupModels)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting lineups by match ID: %w", result.Error)
	}

	return r.mapper.ModelListToDomain(lineupModels), nil
}

func (r *LineupRepositoryImpl) GetLineupsByPlayerID(ctx context.Context, playerID uint64) ([]domain.Lineup, error) {
	var lineupModels []model.Lineup
	result := r.db.WithContext(ctx).
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Where("player_id = ?", playerID).
		Find(&lineupModels)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting lineups by player ID: %w", result.Error)
	}

	return r.mapper.ModelListToDomain(lineupModels), nil
}

func (r *LineupRepositoryImpl) GetPaginatedLineups(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Lineup, int64, error) {
	var lineupModels []model.Lineup
	var total int64

	// Count total records
	countQuery := r.db.WithContext(ctx).Model(&model.Lineup{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total lineups: %w", err)
	}

	// Build the data query with eager loading
	query := r.db.WithContext(ctx).Model(&model.Lineup{}).
		Preload("Player").
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam")

	// Apply sorting if provided
	if sort != "" && (order == "asc" || order == "desc") {
		query = query.Order(fmt.Sprintf("`%s` %s", sort, order))
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&lineupModels).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching lineups: %w", err)
	}

	return r.mapper.ModelListToDomain(lineupModels), total, nil
}

func (r *LineupRepositoryImpl) UpdateLineup(ctx context.Context, id uint64, lineup *domain.Lineup) error {
	lineupModel := r.mapper.DomainToModel(lineup)
	return r.db.WithContext(ctx).
		Model(&model.Lineup{}).
		Where("id = ?", id).
		Select("*").
		Updates(lineupModel).Error
}

func (r *LineupRepositoryImpl) DeleteLineup(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Lineup{}, "id = ?", id).Error
}

func (r *LineupRepositoryImpl) GetStartingLineupsByMatchID(ctx context.Context, matchID uint64) ([]domain.Lineup, error) {
	var lineupModels []model.Lineup
	result := r.db.WithContext(ctx).
		Preload("Player").
		Where("match_id = ? AND starting = ?", matchID, true).
		Find(&lineupModels)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting starting lineups by match ID: %w", result.Error)
	}

	return r.mapper.ModelListToDomain(lineupModels), nil
}

func (r *LineupRepositoryImpl) GetSubstitutesLineupsByMatchID(ctx context.Context, matchID uint64) ([]domain.Lineup, error) {
	var lineupModels []model.Lineup
	result := r.db.WithContext(ctx).
		Preload("Player").
		Where("match_id = ? AND starting = ?", matchID, false).
		Find(&lineupModels)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting substitute lineups by match ID: %w", result.Error)
	}

	return r.mapper.ModelListToDomain(lineupModels), nil
}

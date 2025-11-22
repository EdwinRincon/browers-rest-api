package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/adapter/persistence"
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
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
		Preload(constants.PreloadMatchHomeTeam).
		Preload(constants.PreloadMatchAwayTeam).
		Where(constants.QueryIDEquals, id).
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
		Preload(constants.PreloadMatchHomeTeam).
		Preload(constants.PreloadMatchAwayTeam).
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
		Preload(constants.PreloadMatchHomeTeam).
		Preload(constants.PreloadMatchAwayTeam)

	// Apply sorting (safe and validated)
	col, raw, err := BuildOrderClause(EntityLineup, sort, order)
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
	if err := query.Find(&lineupModels).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching lineups: %w", err)
	}

	return r.mapper.ModelListToDomain(lineupModels), total, nil
}

func (r *LineupRepositoryImpl) UpdateLineup(ctx context.Context, id uint64, lineup *domain.Lineup) error {
	lineupModel := r.mapper.DomainToModel(lineup)
	return r.db.WithContext(ctx).
		Model(&model.Lineup{}).
		Where(constants.QueryIDEquals, id).
		Select("*").
		Updates(lineupModel).Error
}

func (r *LineupRepositoryImpl) DeleteLineup(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Lineup{}, constants.QueryIDEquals, id).Error
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

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

type TeamStatsRepositoryImpl struct {
	db     *gorm.DB
	mapper *persistence.TeamStatsPersistenceMapper
}

func NewTeamStatsRepository(db *gorm.DB) *TeamStatsRepositoryImpl {
	return &TeamStatsRepositoryImpl{
		db:     db,
		mapper: persistence.NewTeamStatsPersistenceMapper(),
	}
}

func (tsr *TeamStatsRepositoryImpl) CreateTeamStats(ctx context.Context, teamStats *domain.TeamStats) error {
	modelTeamStats := tsr.mapper.DomainToModel(teamStats)
	return tsr.db.WithContext(ctx).Create(modelTeamStats).Error
}

func (tsr *TeamStatsRepositoryImpl) GetTeamStatsByID(ctx context.Context, id uint64) (*domain.TeamStats, error) {
	var teamStats model.TeamStat
	result := tsr.db.WithContext(ctx).
		Preload("Team").
		Preload("Season").
		Where("id = ?", id).
		First(&teamStats)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return tsr.mapper.ModelToDomain(&teamStats), nil
}

func (tsr *TeamStatsRepositoryImpl) GetTeamStatsBySeasonAndTeam(ctx context.Context, seasonID, teamID uint64) (*domain.TeamStats, error) {
	var teamStats model.TeamStat
	result := tsr.db.WithContext(ctx).
		Preload("Team").
		Preload("Season").
		Where("season_id = ? AND team_id = ?", seasonID, teamID).
		First(&teamStats)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting team stats by season and team: %w", result.Error)
	}
	return tsr.mapper.ModelToDomain(&teamStats), nil
}

func (tsr *TeamStatsRepositoryImpl) GetTeamStatsBySeasonID(ctx context.Context, seasonID uint64) ([]domain.TeamStats, error) {
	var teamStats []model.TeamStat
	result := tsr.db.WithContext(ctx).
		Preload("Team").
		Preload("Season").
		Where("season_id = ?", seasonID).
		Order("`rank` ASC").
		Find(&teamStats)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting team stats by season: %w", result.Error)
	}
	return tsr.mapper.ModelListToDomain(teamStats), nil
}

func (tsr *TeamStatsRepositoryImpl) GetTeamStatsByTeamID(ctx context.Context, teamID uint64) ([]domain.TeamStats, error) {
	var teamStats []model.TeamStat
	result := tsr.db.WithContext(ctx).
		Preload("Team").
		Preload("Season").
		Where("team_id = ?", teamID).
		Order("created_at DESC").
		Find(&teamStats)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting team stats by team: %w", result.Error)
	}
	return tsr.mapper.ModelListToDomain(teamStats), nil
}

// GetPaginatedTeamStats retrieves a paginated list of team stats with their relationships and total count.
func (tsr *TeamStatsRepositoryImpl) GetPaginatedTeamStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.TeamStats, int64, error) {
	var teamStats []model.TeamStat
	var total int64

	// Count total records
	countQuery := tsr.db.WithContext(ctx).Model(&model.TeamStat{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total team stats: %w", err)
	}

	query := tsr.db.WithContext(ctx).Model(&model.TeamStat{}).
		Preload("Team").
		Preload("Season")

	// Apply sorting (safe and validated)
	col, raw, err := BuildOrderClause(EntityTeamStats, sort, order)
	if err != nil {
		return nil, 0, fmt.Errorf("error building sort clause: %w", err)
	}

	if raw != "" {
		query = query.Order(raw)
	} else {
		query = query.Order(col)
	}

	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&teamStats).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching team stats: %w", err)
	}

	return tsr.mapper.ModelListToDomain(teamStats), total, nil
}

func (tsr *TeamStatsRepositoryImpl) UpdateTeamStats(ctx context.Context, id uint64, teamStats *domain.TeamStats) error {
	modelTeamStats := tsr.mapper.DomainToModel(teamStats)
	return tsr.db.WithContext(ctx).
		Model(&model.TeamStat{}).
		Where("id = ?", id).
		Select("*").
		Updates(modelTeamStats).Error
}

func (tsr *TeamStatsRepositoryImpl) DeleteTeamStats(ctx context.Context, id uint64) error {
	return tsr.db.WithContext(ctx).Delete(&model.TeamStat{}, "id = ?", id).Error
}

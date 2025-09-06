package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type TeamStatsRepository interface {
	CreateTeamStats(ctx context.Context, teamStats *model.TeamStat) error
	GetTeamStatsByID(ctx context.Context, id uint64) (*model.TeamStat, error)
	GetTeamStatsBySeasonAndTeam(ctx context.Context, seasonID, teamID uint64) (*model.TeamStat, error)
	GetTeamStatsBySeasonID(ctx context.Context, seasonID uint64) ([]model.TeamStat, error)
	GetTeamStatsByTeamID(ctx context.Context, teamID uint64) ([]model.TeamStat, error)
	GetPaginatedTeamStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.TeamStat, int64, error)
	UpdateTeamStats(ctx context.Context, id uint64, teamStats *model.TeamStat) error
	DeleteTeamStats(ctx context.Context, id uint64) error
}

type TeamStatsRepositoryImpl struct {
	db *gorm.DB
}

func NewTeamStatsRepository(db *gorm.DB) TeamStatsRepository {
	return &TeamStatsRepositoryImpl{db: db}
}

func (tsr *TeamStatsRepositoryImpl) CreateTeamStats(ctx context.Context, teamStats *model.TeamStat) error {
	return tsr.db.WithContext(ctx).Create(teamStats).Error
}

func (tsr *TeamStatsRepositoryImpl) GetTeamStatsByID(ctx context.Context, id uint64) (*model.TeamStat, error) {
	var teamStats model.TeamStat
	result := tsr.db.WithContext(ctx).
		Preload("Team").
		Preload("Season").
		Where("id = ?", id).
		First(&teamStats)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &teamStats, result.Error
}

func (tsr *TeamStatsRepositoryImpl) GetTeamStatsBySeasonAndTeam(ctx context.Context, seasonID, teamID uint64) (*model.TeamStat, error) {
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
	return &teamStats, nil
}

func (tsr *TeamStatsRepositoryImpl) GetTeamStatsBySeasonID(ctx context.Context, seasonID uint64) ([]model.TeamStat, error) {
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
	return teamStats, nil
}

func (tsr *TeamStatsRepositoryImpl) GetTeamStatsByTeamID(ctx context.Context, teamID uint64) ([]model.TeamStat, error) {
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
	return teamStats, nil
}

// GetPaginatedTeamStats retrieves a paginated list of team stats with their relationships and total count.
func (tsr *TeamStatsRepositoryImpl) GetPaginatedTeamStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.TeamStat, int64, error) {
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

	if sort != "" && (order == "asc" || order == "desc") {
		// Escape the sort field with backticks to handle reserved words
		query = query.Order(fmt.Sprintf("`%s` %s", sort, order))
	}

	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&teamStats).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching team stats: %w", err)
	}

	return teamStats, total, nil
}

func (tsr *TeamStatsRepositoryImpl) UpdateTeamStats(ctx context.Context, id uint64, teamStats *model.TeamStat) error {
	return tsr.db.WithContext(ctx).
		Model(&model.TeamStat{}).
		Where("id = ?", id).
		Select("*").
		Updates(teamStats).Error
}

func (tsr *TeamStatsRepositoryImpl) DeleteTeamStats(ctx context.Context, id uint64) error {
	return tsr.db.WithContext(ctx).Delete(&model.TeamStat{}, "id = ?", id).Error
}

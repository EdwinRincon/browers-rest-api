package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EdwinRincon/browersfc-api/adapter/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
	"gorm.io/gorm"
)

const wherePlayerTeamIDs = "player_id = ? AND team_id = ? AND season_id = ?"

// PlayerTeamRepositoryImpl implements domain.PlayerTeamRepository interface.
type PlayerTeamRepositoryImpl struct {
	db     *gorm.DB
	mapper *mapper.PlayerTeamMapper
}

// NewPlayerTeamRepository creates a new PlayerTeamRepository implementation.
func NewPlayerTeamRepository(db *gorm.DB) *PlayerTeamRepositoryImpl {
	return &PlayerTeamRepositoryImpl{
		db:     db,
		mapper: mapper.NewPlayerTeamMapper(),
	}
}

// Create adds a new player-team association
func (r *PlayerTeamRepositoryImpl) Create(ctx context.Context, playerTeam *domain.PlayerTeam) error {
	modelPlayerTeam := r.mapper.DomainToModel(playerTeam)
	if err := r.db.WithContext(ctx).Create(modelPlayerTeam).Error; err != nil {
		return fmt.Errorf("failed to create player team association: %w", err)
	}

	*playerTeam = *r.mapper.ModelToDomain(modelPlayerTeam)
	return nil
}

// GetByPlayerID gets all team associations for a player
func (r *PlayerTeamRepositoryImpl) GetByPlayerID(ctx context.Context, playerID uint64) ([]domain.PlayerTeam, error) {
	var playerTeams []model.PlayerTeam

	result := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Team").
		Preload("Season").
		Where("player_id = ?", playerID).
		Find(&playerTeams)

	if result.Error != nil {
		return nil, result.Error
	}

	return r.mapper.ModelListToDomain(playerTeams), nil
}

// DeleteByPlayerID removes all team associations for a player
func (r *PlayerTeamRepositoryImpl) DeleteByPlayerID(ctx context.Context, playerID uint64) error {
	if err := r.db.WithContext(ctx).
		Where("player_id = ?", playerID).
		Delete(&model.PlayerTeam{}).Error; err != nil {
		return fmt.Errorf("failed to delete player team associations: %w", err)
	}
	return nil
}

// GetPlayerTeamByID retrieves a player-team relationship by its ID
func (r *PlayerTeamRepositoryImpl) GetPlayerTeamByID(ctx context.Context, id uint64) (*domain.PlayerTeam, error) {
	var playerTeam model.PlayerTeam
	result := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Team").
		Preload("Season").
		First(&playerTeam, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting player team relationship by ID: %w", result.Error)
	}
	return r.mapper.ModelToDomain(&playerTeam), nil
}

// GetPlayerTeamsByTeamID retrieves all player relationships for a specific team
func (r *PlayerTeamRepositoryImpl) GetPlayerTeamsByTeamID(ctx context.Context, teamID uint64) ([]domain.PlayerTeam, error) {
	var playerTeams []model.PlayerTeam
	err := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Team").
		Preload("Season").
		Where("team_id = ?", teamID).
		Find(&playerTeams).Error

	if err != nil {
		return nil, fmt.Errorf("error getting player teams by team ID: %w", err)
	}
	return r.mapper.ModelListToDomain(playerTeams), nil
}

// GetPlayerTeamsBySeasonID retrieves all player-team relationships for a specific season
func (r *PlayerTeamRepositoryImpl) GetPlayerTeamsBySeasonID(ctx context.Context, seasonID uint64) ([]domain.PlayerTeam, error) {
	var playerTeams []model.PlayerTeam
	err := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Team").
		Preload("Season").
		Where("season_id = ?", seasonID).
		Find(&playerTeams).Error

	if err != nil {
		return nil, fmt.Errorf("error getting player teams by season ID: %w", err)
	}
	return r.mapper.ModelListToDomain(playerTeams), nil
}

// GetPaginatedPlayerTeams retrieves a paginated list of player-team relationships
func (r *PlayerTeamRepositoryImpl) GetPaginatedPlayerTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.PlayerTeam, int64, error) {
	var playerTeams []model.PlayerTeam
	var total int64

	// Count total records
	countQuery := r.db.WithContext(ctx).Model(&model.PlayerTeam{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total player teams: %w", err)
	}

	// Build the data query with eager loading
	query := r.db.WithContext(ctx).Model(&model.PlayerTeam{}).
		Preload("Player").
		Preload("Team").
		Preload("Season")

	// Apply sorting if provided
	if sort != "" && (order == "asc" || order == "desc") {
		// Escape the sort field with backticks to handle reserved words
		query = query.Order(fmt.Sprintf("`%s` %s", sort, order))
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&playerTeams).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching player teams: %w", err)
	}

	return r.mapper.ModelListToDomain(playerTeams), total, nil
}

// UpdatePlayerTeam updates an existing player-team relationship
func (r *PlayerTeamRepositoryImpl) UpdatePlayerTeam(ctx context.Context, playerTeam *domain.PlayerTeam) error {
	modelPlayerTeam := r.mapper.DomainToModel(playerTeam)
	return r.db.WithContext(ctx).
		Model(&model.PlayerTeam{}).
		Where("id = ?", modelPlayerTeam.ID).
		Updates(map[string]interface{}{
			"player_id":  modelPlayerTeam.PlayerID,
			"team_id":    modelPlayerTeam.TeamID,
			"season_id":  modelPlayerTeam.SeasonID,
			"start_date": modelPlayerTeam.StartDate,
			"end_date":   modelPlayerTeam.EndDate,
		}).Error
}

// DeletePlayerTeam soft-deletes a player-team relationship
func (r *PlayerTeamRepositoryImpl) DeletePlayerTeam(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.PlayerTeam{}).Error
}

// checkExistingDates checks for overlapping dates for a specific query
func (r *PlayerTeamRepositoryImpl) checkExistingDates(query *gorm.DB, startDate time.Time, endDate *time.Time) (bool, error) {
	var exists []int

	// overlap detection
	endDateValue := "9999-12-31"
	if endDate != nil {
		endDateValue = endDate.Format("2006-01-02 15:04:05")
	}

	overlapQuery := query.Session(&gorm.Session{}).
		Where("start_date <= ? AND COALESCE(end_date, '9999-12-31') >= ?",
			endDateValue, startDate)

	err := overlapQuery.Select("1").Limit(1).Find(&exists).Error
	if err != nil {
		return false, fmt.Errorf("error checking date overlaps: %w", err)
	}

	return len(exists) > 0, nil
}

// CheckOverlappingDates checks if there are overlapping dates for the same player-team-season combination
func (r *PlayerTeamRepositoryImpl) CheckOverlappingDates(ctx context.Context, data domain.OverlapCheckData) (bool, error) {
	return r.checkDateOverlaps(ctx, data)
}

// checkDateOverlaps is a helper function to reduce parameter count
func (r *PlayerTeamRepositoryImpl) checkDateOverlaps(ctx context.Context, data domain.OverlapCheckData) (bool, error) {
	// Build the base query to find records with the same player-team-season
	query := r.db.WithContext(ctx).Model(&model.PlayerTeam{}).
		Where(wherePlayerTeamIDs, data.PlayerID, data.TeamID, data.SeasonID)

	// If updating an existing record, exclude it from the check
	if data.IsUpdate && data.ID > 0 {
		query = query.Where("id != ?", data.ID)
	}

	// Perform the date overlap checks
	return r.checkExistingDates(query, data.StartDate, data.EndDate)
}

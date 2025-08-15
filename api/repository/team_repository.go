package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeamByID(ctx context.Context, id uint64) (*model.Team, error)
	GetActiveTeamByName(ctx context.Context, fullName string) (*model.Team, error)
	GetUnscopedTeamByName(ctx context.Context, fullName string) (*model.Team, error)
	GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Team, int64, error)
	UpdateTeam(ctx context.Context, team *model.Team) error
	DeleteTeam(ctx context.Context, id uint64) error
	RestoreAndUpdateTeam(ctx context.Context, team *model.Team) error
}

type TeamRepositoryImpl struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &TeamRepositoryImpl{db: db}
}

func (tr *TeamRepositoryImpl) CreateTeam(ctx context.Context, team *model.Team) error {
	return tr.db.WithContext(ctx).Create(team).Error
}

func (tr *TeamRepositoryImpl) GetTeamByID(ctx context.Context, id uint64) (*model.Team, error) {
	var team model.Team
	err := tr.db.WithContext(ctx).First(&team, id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// GetActiveTeamByName retrieves an active (not deleted) team by its full name.
func (tr *TeamRepositoryImpl) GetActiveTeamByName(ctx context.Context, fullName string) (*model.Team, error) {
	var team model.Team
	result := tr.db.WithContext(ctx).
		Where("short_name = ?", fullName).
		First(&team)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting team by name: %w", result.Error)
	}
	return &team, nil
}

// GetUnscopedTeamByName retrieves a team by its full name, including soft-deleted records.
func (tr *TeamRepositoryImpl) GetUnscopedTeamByName(ctx context.Context, fullName string) (*model.Team, error) {
	var team model.Team
	result := tr.db.WithContext(ctx).
		Unscoped().
		Where("short_name = ?", fullName).
		First(&team)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &team, result.Error
}

// GetPaginatedTeams retrieves a paginated list of teams with total count.
func (tr *TeamRepositoryImpl) GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Team, int64, error) {
	var teams []model.Team
	var total int64

	// Count total records
	countQuery := tr.db.WithContext(ctx).Model(&model.Team{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total teams: %w", err)
	}

	// Build the data query
	query := tr.db.WithContext(ctx).Model(&model.Team{})

	// Apply sorting if provided
	if sort != "" && (order == "asc" || order == "desc") {
		query = query.Order(fmt.Sprintf("%s %s", sort, order))
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&teams).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching teams: %w", err)
	}

	return teams, total, nil
}

func (tr *TeamRepositoryImpl) UpdateTeam(ctx context.Context, team *model.Team) error {
	result := tr.db.WithContext(ctx).Save(team)
	return result.Error
}

func (tr *TeamRepositoryImpl) DeleteTeam(ctx context.Context, id uint64) error {
	result := tr.db.WithContext(ctx).Delete(&model.Team{}, id)
	return result.Error
}

// RestoreAndUpdateTeam restores a soft-deleted team and updates its information in a single transaction
func (tr *TeamRepositoryImpl) RestoreAndUpdateTeam(ctx context.Context, team *model.Team) error {
	return tr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First verify the team exists and is soft-deleted
		var existingTeam model.Team
		if err := tx.Unscoped().Where("id = ?", team.ID).First(&existingTeam).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("team not found: %w", err)
			}
			return fmt.Errorf("failed to find team: %w", err)
		}

		if !existingTeam.DeletedAt.Valid {
			return fmt.Errorf("team is not soft-deleted")
		}

		// Preserve important metadata
		team.CreatedAt = existingTeam.CreatedAt
		team.DeletedAt = gorm.DeletedAt{} // Explicitly set to zero value to restore

		// Save the entire team object (this will update all fields including deleted_at)
		if err := tx.Unscoped().Save(team).Error; err != nil {
			return fmt.Errorf("failed to restore and update team: %w", err)
		}

		// Verify the restoration was successful
		var restoredTeam model.Team
		if err := tx.Where("id = ?", team.ID).First(&restoredTeam).Error; err != nil {
			return fmt.Errorf("failed to verify team restoration: %w", err)
		}

		if restoredTeam.DeletedAt.Valid {
			return fmt.Errorf("team restoration failed: deleted_at field is still set")
		}

		return nil
	})
}

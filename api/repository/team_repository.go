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
	GetTeamByName(ctx context.Context, fullName string) (*model.Team, error)
	GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Team, int64, error)
	UpdateTeam(ctx context.Context, team *model.Team) error
	DeleteTeam(ctx context.Context, id uint64) error
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

func (tr *TeamRepositoryImpl) GetTeamByName(ctx context.Context, fullName string) (*model.Team, error) {
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

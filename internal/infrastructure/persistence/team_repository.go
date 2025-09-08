package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/adapter/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
	"gorm.io/gorm"
)

// TeamRepositoryImpl implements domain.TeamRepository interface.
type TeamRepositoryImpl struct {
	db     *gorm.DB
	mapper *mapper.TeamMapper
}

func NewTeamRepository(db *gorm.DB) *TeamRepositoryImpl {
	return &TeamRepositoryImpl{
		db:     db,
		mapper: mapper.NewTeamMapper(),
	}
}

func (tr *TeamRepositoryImpl) CreateTeam(ctx context.Context, team *domain.Team) error {
	modelTeam := tr.mapper.DomainToModel(team)
	if err := tr.db.WithContext(ctx).Create(modelTeam).Error; err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	// Update domain entity with generated ID and timestamps
	*team = *tr.mapper.ModelToDomain(modelTeam)
	return nil
}

func (tr *TeamRepositoryImpl) GetTeamByID(ctx context.Context, id uint64) (*domain.Team, error) {
	var team model.Team
	result := tr.db.WithContext(ctx).First(&team, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting team by ID: %w", result.Error)
	}

	return tr.mapper.ModelToDomain(&team), nil
}

func (tr *TeamRepositoryImpl) GetTeamByName(ctx context.Context, fullName string) (*domain.Team, error) {
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

	return tr.mapper.ModelToDomain(&team), nil
}

// GetPaginatedTeams retrieves a paginated list of teams with total count.
func (tr *TeamRepositoryImpl) GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Team, int64, error) {
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
		query = query.Order(fmt.Sprintf("`%s` %s", sort, order))
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&teams).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching teams: %w", err)
	}

	return tr.mapper.ModelListToDomain(teams), total, nil
}

func (tr *TeamRepositoryImpl) UpdateTeam(ctx context.Context, team *domain.Team) error {
	modelTeam := tr.mapper.DomainToModel(team)
	result := tr.db.WithContext(ctx).Save(modelTeam)
	if result.Error != nil {
		return fmt.Errorf("failed to update team: %w", result.Error)
	}

	// Update domain entity with new timestamps
	*team = *tr.mapper.ModelToDomain(modelTeam)
	return nil
}

func (tr *TeamRepositoryImpl) DeleteTeam(ctx context.Context, id uint64) error {
	result := tr.db.WithContext(ctx).Delete(&model.Team{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete team: %w", result.Error)
	}
	return nil
}

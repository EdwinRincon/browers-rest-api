package repository

import (
	"context"
	"errors"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

var ErrTeamNotFound = errors.New("team not found")

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *model.Teams) error
	GetTeamByID(ctx context.Context, id uint64) (*model.Teams, error)
	ListTeams(ctx context.Context, page uint64) ([]*model.Teams, error)
	UpdateTeam(ctx context.Context, team *model.Teams) error
	DeleteTeam(ctx context.Context, id uint64) error
}

type TeamRepositoryImpl struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &TeamRepositoryImpl{db: db}
}

func (tr *TeamRepositoryImpl) CreateTeam(ctx context.Context, team *model.Teams) error {
	return tr.db.WithContext(ctx).Create(team).Error
}

func (tr *TeamRepositoryImpl) GetTeamByID(ctx context.Context, id uint64) (*model.Teams, error) {
	var team model.Teams
	err := tr.db.WithContext(ctx).First(&team, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}
	return &team, nil
}

func (tr *TeamRepositoryImpl) ListTeams(ctx context.Context, page uint64) ([]*model.Teams, error) {
	var teams []*model.Teams
	offset := (page - 1) * 10
	err := tr.db.WithContext(ctx).Offset(int(offset)).Limit(10).Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (tr *TeamRepositoryImpl) UpdateTeam(ctx context.Context, team *model.Teams) error {
	result := tr.db.WithContext(ctx).Save(team)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTeamNotFound
	}
	return nil
}

func (tr *TeamRepositoryImpl) DeleteTeam(ctx context.Context, id uint64) error {
	result := tr.db.WithContext(ctx).Delete(&model.Teams{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTeamNotFound
	}
	return nil
}
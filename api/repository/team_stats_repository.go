package repository

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type TeamStatsRepository interface {
	CreateTeamStat(ctx context.Context, teamStat *model.TeamStat) error
	GetTeamStatByID(ctx context.Context, id uint64) (*model.TeamStat, error)
	ListTeamStats(ctx context.Context, page uint64) ([]*model.TeamStat, error)
	UpdateTeamStat(ctx context.Context, teamStat *model.TeamStat) error
	DeleteTeamStat(ctx context.Context, id uint64) error
}

type TeamStatsRepositoryImpl struct {
	db *gorm.DB
}

func NewTeamStatsRepository(db *gorm.DB) TeamStatsRepository {
	return &TeamStatsRepositoryImpl{db: db}
}

func (r *TeamStatsRepositoryImpl) CreateTeamStat(ctx context.Context, teamStat *model.TeamStat) error {
	return r.db.WithContext(ctx).Create(teamStat).Error
}

func (r *TeamStatsRepositoryImpl) GetTeamStatByID(ctx context.Context, id uint64) (*model.TeamStat, error) {
	var teamStat model.TeamStat
	err := r.db.WithContext(ctx).First(&teamStat, id).Error
	if err != nil {
		return nil, err
	}
	return &teamStat, nil
}

func (r *TeamStatsRepositoryImpl) ListTeamStats(ctx context.Context, page uint64) ([]*model.TeamStat, error) {
	var teamStats []*model.TeamStat
	offset := (page - 1) * 10
	err := r.db.WithContext(ctx).Offset(int(offset)).Limit(10).Find(&teamStats).Error
	if err != nil {
		return nil, err
	}
	return teamStats, nil
}

func (r *TeamStatsRepositoryImpl) UpdateTeamStat(ctx context.Context, teamStat *model.TeamStat) error {
	result := r.db.WithContext(ctx).Save(teamStat)
	return result.Error
}

func (r *TeamStatsRepositoryImpl) DeleteTeamStat(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.TeamStat{}, id)
	return result.Error
}

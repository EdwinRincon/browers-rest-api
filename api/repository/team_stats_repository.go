package repository

import (
	"context"
	"errors"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

var ErrTeamStatsNotFound = errors.New("team stats not found")

type TeamStatsRepository interface {
	CreateTeamStats(ctx context.Context, teamStats *model.TeamStat) error
	GetTeamStatsByID(ctx context.Context, id uint64) (*model.TeamStat, error)
	ListTeamStats(ctx context.Context, page uint64) ([]*model.TeamStat, error)
	UpdateTeamStats(ctx context.Context, teamStats *model.TeamStat) error
	DeleteTeamStats(ctx context.Context, id uint64) error
}

type TeamStatsRepositoryImpl struct {
	db *gorm.DB
}

func NewTeamStatsRepository(db *gorm.DB) TeamStatsRepository {
	return &TeamStatsRepositoryImpl{db: db}
}

func (r *TeamStatsRepositoryImpl) CreateTeamStats(ctx context.Context, teamStats *model.TeamStat) error {
	return r.db.WithContext(ctx).Create(teamStats).Error
}

func (r *TeamStatsRepositoryImpl) GetTeamStatsByID(ctx context.Context, id uint64) (*model.TeamStat, error) {
	var teamStats model.TeamStat
	err := r.db.WithContext(ctx).First(&teamStats, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamStatsNotFound
		}
		return nil, err
	}
	return &teamStats, nil
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

func (r *TeamStatsRepositoryImpl) UpdateTeamStats(ctx context.Context, teamStats *model.TeamStat) error {
	result := r.db.WithContext(ctx).Save(teamStats)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTeamStatsNotFound
	}
	return nil
}

func (r *TeamStatsRepositoryImpl) DeleteTeamStats(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.TeamStat{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTeamStatsNotFound
	}
	return nil
}

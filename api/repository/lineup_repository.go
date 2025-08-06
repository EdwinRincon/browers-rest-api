package repository

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type LineupRepository interface {
	CreateLineup(ctx context.Context, lineup *model.Lineup) error
	GetLineupByID(ctx context.Context, id uint64) (*model.Lineup, error)
	ListLineups(ctx context.Context, page uint64) ([]*model.Lineup, error)
	UpdateLineup(ctx context.Context, lineup *model.Lineup) error
	DeleteLineup(ctx context.Context, id uint64) error
	GetLineupsByMatch(ctx context.Context, matchID uint64) ([]*model.Lineup, error)
}

type LineupRepositoryImpl struct {
	db *gorm.DB
}

func NewLineupRepository(db *gorm.DB) LineupRepository {
	return &LineupRepositoryImpl{db: db}
}

func (r *LineupRepositoryImpl) CreateLineup(ctx context.Context, lineup *model.Lineup) error {
	return r.db.WithContext(ctx).Create(lineup).Error
}

func (r *LineupRepositoryImpl) GetLineupByID(ctx context.Context, id uint64) (*model.Lineup, error) {
	var lineup model.Lineup
	err := r.db.WithContext(ctx).Preload("Player").Preload("Match").First(&lineup, id).Error
	if err != nil {
		return nil, err
	}
	return &lineup, nil
}

func (r *LineupRepositoryImpl) ListLineups(ctx context.Context, page uint64) ([]*model.Lineup, error) {
	var lineups []*model.Lineup
	const itemsPerPage = 10
	offset := (page - 1) * itemsPerPage

	err := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Match").
		Order("id ASC").
		Limit(int(itemsPerPage)).
		Offset(int(offset)).
		Find(&lineups).Error

	if err != nil {
		return nil, err
	}

	return lineups, nil
}

func (r *LineupRepositoryImpl) UpdateLineup(ctx context.Context, lineup *model.Lineup) error {
	result := r.db.WithContext(ctx).Save(lineup)
	return result.Error
}

func (r *LineupRepositoryImpl) DeleteLineup(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.Lineup{}, id)
	return result.Error
}

func (r *LineupRepositoryImpl) GetLineupsByMatch(ctx context.Context, matchID uint64) ([]*model.Lineup, error) {
	var lineups []*model.Lineup
	err := r.db.WithContext(ctx).Preload("Player").Where("match_id = ?", matchID).Find(&lineups).Error
	if err != nil {
		return nil, err
	}
	return lineups, nil
}

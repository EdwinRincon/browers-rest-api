package repository

import (
	"context"
	"errors"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type MatchRepository interface {
	CreateMatch(ctx context.Context, match *model.Match) error
	GetMatchByID(ctx context.Context, id uint64) (*model.Match, error)
	ListMatches(ctx context.Context, page, pageSize uint64) ([]*model.Match, error)
	UpdateMatch(ctx context.Context, match *model.Match) error
	DeleteMatch(ctx context.Context, id uint64) error
}

type MatchRepositoryImpl struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &MatchRepositoryImpl{db: db}
}

func (mr *MatchRepositoryImpl) CreateMatch(ctx context.Context, match *model.Match) error {
	return mr.db.WithContext(ctx).Create(match).Error
}

func (mr *MatchRepositoryImpl) GetMatchByID(ctx context.Context, id uint64) (*model.Match, error) {
	var match model.Match
	err := mr.db.WithContext(ctx).Preload("Season").Preload("MVPPlayer").First(&match, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrMatchNotFound
		}
		return nil, err
	}
	return &match, nil
}

func (mr *MatchRepositoryImpl) ListMatches(ctx context.Context, page, pageSize uint64) ([]*model.Match, error) {
	var matches []*model.Match
	offset := (page - 1) * pageSize
	err := mr.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Lineups").
		Preload("Season").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (mr *MatchRepositoryImpl) GetNextScheduledMatch(ctx context.Context) (*model.Match, error) {
	var match model.Match
	err := mr.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Where("date >= CURRENT_DATE").
		Where("status = ?", "scheduled").
		Where("home_team_id = ?", 1).Or("away_team_id =?", 1).First(&match).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrMatchNotFound
		}
		return nil, err
	}
	return &match, nil
}

func (mr *MatchRepositoryImpl) UpdateMatch(ctx context.Context, match *model.Match) error {
	result := mr.db.WithContext(ctx).Save(match)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return constants.ErrMatchNotFound
	}
	return nil
}

func (mr *MatchRepositoryImpl) DeleteMatch(ctx context.Context, id uint64) error {
	result := mr.db.WithContext(ctx).Delete(&model.Match{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return constants.ErrMatchNotFound
	}
	return nil
}

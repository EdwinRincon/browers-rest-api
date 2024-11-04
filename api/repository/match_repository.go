package repository

import (
	"context"
	"errors"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

var ErrMatchNotFound = errors.New("match not found")

type MatchRepository interface {
	CreateMatch(ctx context.Context, match *model.Matches) error
	GetMatchByID(ctx context.Context, id uint64) (*model.Matches, error)
	ListMatches(ctx context.Context, page, pageSize uint64) ([]*model.Matches, error)
	UpdateMatch(ctx context.Context, match *model.Matches) error
	DeleteMatch(ctx context.Context, id uint64) error
}

type MatchRepositoryImpl struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &MatchRepositoryImpl{db: db}
}

func (mr *MatchRepositoryImpl) CreateMatch(ctx context.Context, match *model.Matches) error {
	return mr.db.WithContext(ctx).Create(match).Error
}

func (mr *MatchRepositoryImpl) GetMatchByID(ctx context.Context, id uint64) (*model.Matches, error) {
	var match model.Matches
	err := mr.db.WithContext(ctx).Preload("Seasons").Preload("MVPPlayer").First(&match, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}
	return &match, nil
}

func (mr *MatchRepositoryImpl) ListMatches(ctx context.Context, page, pageSize uint64) ([]*model.Matches, error) {
	var matches []*model.Matches
	offset := (page - 1) * pageSize
	err := mr.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Lineups").
		Preload("Seasons").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (mr *MatchRepositoryImpl) UpdateMatch(ctx context.Context, match *model.Matches) error {
	result := mr.db.WithContext(ctx).Save(match)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrMatchNotFound
	}
	return nil
}

func (mr *MatchRepositoryImpl) DeleteMatch(ctx context.Context, id uint64) error {
	result := mr.db.WithContext(ctx).Delete(&model.Matches{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrMatchNotFound
	}
	return nil
}

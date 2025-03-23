package repository

import (
	"context"
	"errors"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

var ErrSeasonNotFound = errors.New("season not found")

type SeasonRepository interface {
	CreateSeason(ctx context.Context, season *model.Season) error
	GetSeasonByID(ctx context.Context, id uint8) (*model.Season, error)
	GetAllSeasons(ctx context.Context, page uint64) ([]*model.Season, error)
	UpdateSeason(ctx context.Context, season *model.Season) error
	DeleteSeason(ctx context.Context, id uint8) error
}

type SeasonRepositoryImpl struct {
	db *gorm.DB
}

func NewSeasonRepository(db *gorm.DB) SeasonRepository {
	return &SeasonRepositoryImpl{db: db}
}

func (sr *SeasonRepositoryImpl) CreateSeason(ctx context.Context, season *model.Season) error {
	return sr.db.WithContext(ctx).Create(season).Error
}

func (sr *SeasonRepositoryImpl) GetSeasonByID(ctx context.Context, id uint8) (*model.Season, error) {
	var season model.Season
	err := sr.db.WithContext(ctx).Preload("Match").Preload("Article").Preload("TeamStat").First(&season, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSeasonNotFound
		}
		return nil, err
	}
	return &season, nil
}

func (sr *SeasonRepositoryImpl) GetAllSeasons(ctx context.Context, page uint64) ([]*model.Season, error) {
	var seasons []*model.Season
	offset := (page - 1) * 10
	err := sr.db.WithContext(ctx).Offset(int(offset)).Limit(10).Find(&seasons).Error
	if err != nil {
		return nil, err
	}
	return seasons, nil
}

func (sr *SeasonRepositoryImpl) UpdateSeason(ctx context.Context, season *model.Season) error {
	result := sr.db.WithContext(ctx).Save(season)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSeasonNotFound
	}
	return nil
}

func (sr *SeasonRepositoryImpl) DeleteSeason(ctx context.Context, id uint8) error {
	result := sr.db.WithContext(ctx).Delete(&model.Season{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSeasonNotFound
	}
	return nil
}

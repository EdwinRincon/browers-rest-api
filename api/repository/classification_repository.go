package repository

import (
	"context"
	"errors"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

var ErrClassificationNotFound = errors.New("classification not found")

type ClassificationRepository interface {
	CreateClassification(ctx context.Context, classification *model.Classifications) error
	GetClassificationByID(ctx context.Context, id uint64) (*model.Classifications, error)
	ListClassifications(ctx context.Context, page uint64) ([]*model.Classifications, error)
	UpdateClassification(ctx context.Context, classification *model.Classifications) error
	DeleteClassification(ctx context.Context, id uint64) error
	GetClassificationBySeason(ctx context.Context, seasonID uint8) ([]*model.Classifications, error)
}

type ClassificationRepositoryImpl struct {
	db *gorm.DB
}

func NewClassificationRepository(db *gorm.DB) ClassificationRepository {
	return &ClassificationRepositoryImpl{db: db}
}

func (r *ClassificationRepositoryImpl) CreateClassification(ctx context.Context, classification *model.Classifications) error {
	return r.db.WithContext(ctx).Create(classification).Error
}

func (r *ClassificationRepositoryImpl) GetClassificationByID(ctx context.Context, id uint64) (*model.Classifications, error) {
	var classification model.Classifications
	err := r.db.WithContext(ctx).Preload("TeamsStats").Preload("Seasons").First(&classification, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClassificationNotFound
		}
		return nil, err
	}
	return &classification, nil
}

func (r *ClassificationRepositoryImpl) ListClassifications(ctx context.Context, page uint64) ([]*model.Classifications, error) {
	var classifications []*model.Classifications
	offset := (page - 1) * 10
	err := r.db.WithContext(ctx).Preload("TeamsStats").Preload("Seasons").Offset(int(offset)).Limit(10).Find(&classifications).Error
	if err != nil {
		return nil, err
	}
	return classifications, nil
}

func (r *ClassificationRepositoryImpl) UpdateClassification(ctx context.Context, classification *model.Classifications) error {
	result := r.db.WithContext(ctx).Save(classification)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrClassificationNotFound
	}
	return nil
}

func (r *ClassificationRepositoryImpl) DeleteClassification(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.Classifications{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrClassificationNotFound
	}
	return nil
}

func (r *ClassificationRepositoryImpl) GetClassificationBySeason(ctx context.Context, seasonID uint8) ([]*model.Classifications, error) {
	var classifications []*model.Classifications
	err := r.db.WithContext(ctx).Preload("TeamsStats").Preload("Seasons").Where("seasons_id = ?", seasonID).Order("points DESC").Find(&classifications).Error
	if err != nil {
		return nil, err
	}
	return classifications, nil
}

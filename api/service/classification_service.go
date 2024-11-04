package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type ClassificationService interface {
	CreateClassification(ctx context.Context, classification *model.Classifications) error
	GetClassificationByID(ctx context.Context, id uint64) (*model.Classifications, error)
	ListClassifications(ctx context.Context, page uint64) ([]*model.Classifications, error)
	UpdateClassification(ctx context.Context, classification *model.Classifications) error
	DeleteClassification(ctx context.Context, id uint64) error
	GetClassificationBySeason(ctx context.Context, seasonID uint8) ([]*model.Classifications, error)
}

type classificationService struct {
	ClassificationRepository repository.ClassificationRepository
}

func NewClassificationService(classificationRepo repository.ClassificationRepository) ClassificationService {
	return &classificationService{
		ClassificationRepository: classificationRepo,
	}
}

func (s *classificationService) CreateClassification(ctx context.Context, classification *model.Classifications) error {
	return s.ClassificationRepository.CreateClassification(ctx, classification)
}

func (s *classificationService) GetClassificationByID(ctx context.Context, id uint64) (*model.Classifications, error) {
	return s.ClassificationRepository.GetClassificationByID(ctx, id)
}

func (s *classificationService) ListClassifications(ctx context.Context, page uint64) ([]*model.Classifications, error) {
	return s.ClassificationRepository.ListClassifications(ctx, page)
}

func (s *classificationService) UpdateClassification(ctx context.Context, classification *model.Classifications) error {
	return s.ClassificationRepository.UpdateClassification(ctx, classification)
}

func (s *classificationService) DeleteClassification(ctx context.Context, id uint64) error {
	return s.ClassificationRepository.DeleteClassification(ctx, id)
}

func (s *classificationService) GetClassificationBySeason(ctx context.Context, seasonID uint8) ([]*model.Classifications, error) {
	return s.ClassificationRepository.GetClassificationBySeason(ctx, seasonID)
}

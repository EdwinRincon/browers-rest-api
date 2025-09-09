package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

// ArticlePersistenceMapper handles persistence layer conversions for Article entity
type ArticlePersistenceMapper struct{}

func NewArticlePersistenceMapper() *ArticlePersistenceMapper {
	return &ArticlePersistenceMapper{}
}

// DomainToModel converts a domain Article to model Article for persistence
func (m *ArticlePersistenceMapper) DomainToModel(entity *domain.Article) *model.Article {
	if entity == nil {
		return nil
	}

	return &model.Article{
		ID:        entity.ID,
		Title:     entity.Title,
		Content:   entity.Content,
		ImgBanner: entity.ImgBanner,
		Date:      entity.Date,
		SeasonID:  entity.SeasonID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ModelToDomain converts a model Article to domain Article for business logic
func (m *ArticlePersistenceMapper) ModelToDomain(model *model.Article) *domain.Article {
	if model == nil {
		return nil
	}

	return &domain.Article{
		ID:        model.ID,
		Title:     model.Title,
		Content:   model.Content,
		ImgBanner: model.ImgBanner,
		Date:      model.Date,
		SeasonID:  model.SeasonID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// ModelListToDomain converts a slice of model Article to domain Article for business logic
func (m *ArticlePersistenceMapper) ModelListToDomain(models []model.Article) []domain.Article {
	if models == nil {
		return nil
	}

	domains := make([]domain.Article, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			domains[i] = *domain
		}
	}

	return domains
}

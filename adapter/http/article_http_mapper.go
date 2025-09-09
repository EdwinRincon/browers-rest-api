package http

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// ArticleHTTPMapper handles HTTP layer conversions for Article entity
type ArticleHTTPMapper struct{}

func NewArticleHTTPMapper() *ArticleHTTPMapper {
	return &ArticleHTTPMapper{}
}

// DTOToDomain converts a CreateArticleRequest DTO to a domain Article entity
func (m *ArticleHTTPMapper) DTOToDomain(dto *dto.CreateArticleRequest) *domain.Article {
	if dto == nil {
		return nil
	}

	return &domain.Article{
		Title:     dto.Title,
		Content:   dto.Content,
		ImgBanner: dto.ImgBanner,
		Date:      dto.Date,
		SeasonID:  dto.SeasonID,
	}
}

// UpdateDTOToDomain converts an UpdateArticleRequest DTO to a domain Article entity
func (m *ArticleHTTPMapper) UpdateDTOToDomain(dto *dto.UpdateArticleRequest, existingArticle *domain.Article) *domain.Article {
	if dto == nil || existingArticle == nil {
		return nil
	}

	updatedArticle := &domain.Article{
		ID:        existingArticle.ID,
		Title:     existingArticle.Title,
		Content:   existingArticle.Content,
		ImgBanner: existingArticle.ImgBanner,
		Date:      existingArticle.Date,
		SeasonID:  existingArticle.SeasonID,
		CreatedAt: existingArticle.CreatedAt,
		UpdatedAt: existingArticle.UpdatedAt,
	}

	if dto.Title != nil {
		updatedArticle.Title = *dto.Title
	}
	if dto.Content != nil {
		updatedArticle.Content = *dto.Content
	}
	if dto.ImgBanner != nil {
		updatedArticle.ImgBanner = *dto.ImgBanner
	}
	if dto.Date != nil {
		updatedArticle.Date = *dto.Date
	}
	if dto.SeasonID != nil {
		updatedArticle.SeasonID = *dto.SeasonID
	}

	return updatedArticle
}

// DomainToDTO converts a domain Article to ArticleResponse DTO
func (m *ArticleHTTPMapper) DomainToDTO(entity *domain.Article) *dto.ArticleResponse {
	if entity == nil {
		return nil
	}

	response := &dto.ArticleResponse{
		ID:        entity.ID,
		Title:     entity.Title,
		Content:   entity.Content,
		ImgBanner: entity.ImgBanner,
		Date:      entity.Date,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Season will be populated by handler if needed

	return response
}

// DomainListToDTO converts a slice of domain Article to ArticleResponse DTOs
func (m *ArticleHTTPMapper) DomainListToDTO(entities []domain.Article) []dto.ArticleResponse {
	if entities == nil {
		return nil
	}

	result := make([]dto.ArticleResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			result[i] = *response
		}
	}
	return result
}

// DomainToShortDTO converts a domain Article to ArticleShort DTO
func (m *ArticleHTTPMapper) DomainToShortDTO(entity *domain.Article) *dto.ArticleShort {
	if entity == nil {
		return nil
	}

	return &dto.ArticleShort{
		ID:    entity.ID,
		Title: entity.Title,
		Date:  entity.Date,
	}
}

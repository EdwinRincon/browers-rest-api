package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// ArticleMapper handles all conversions between DTOs, domain entities, and models for articles.
type ArticleMapper struct{}

func NewArticleMapper() *ArticleMapper {
	return &ArticleMapper{}
}

// DTO to Domain Conversions (HTTP layer)

// DTOToDomain converts a CreateArticleRequest DTO to a domain Article entity.
func (m *ArticleMapper) DTOToDomain(dto *dto.CreateArticleRequest) *domain.Article {
	return &domain.Article{
		Title:     dto.Title,
		Content:   dto.Content,
		ImgBanner: dto.ImgBanner,
		Date:      dto.Date,
		SeasonID:  dto.SeasonID,
	}
}

// UpdateDTOToDomain converts an UpdateArticleRequest DTO to a domain Article entity.
func (m *ArticleMapper) UpdateDTOToDomain(dto *dto.UpdateArticleRequest, existingArticle *domain.Article) *domain.Article {
	updatedArticle := &domain.Article{
		ID:        existingArticle.ID,
		Title:     existingArticle.Title,
		Content:   existingArticle.Content,
		ImgBanner: existingArticle.ImgBanner,
		Date:      existingArticle.Date,
		SeasonID:  existingArticle.SeasonID,
		Season:    existingArticle.Season,
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

// DomainToDTO converts a domain Article entity to an ArticleResponse DTO.
func (m *ArticleMapper) DomainToDTO(article *domain.Article) *dto.ArticleResponse {
	var season dto.SeasonShort
	if article.Season != nil {
		seasonMapper := NewSeasonMapper()
		seasonShort := seasonMapper.DomainToShortDTO(article.Season)
		if seasonShort != nil {
			season = *seasonShort
		}
	}

	return &dto.ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		ImgBanner: article.ImgBanner,
		Date:      article.Date,
		Season:    season,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}
}

// DomainListToDTO converts a slice of domain Article entities to ArticleResponse DTOs.
func (m *ArticleMapper) DomainListToDTO(articles []domain.Article) []dto.ArticleResponse {
	responses := make([]dto.ArticleResponse, len(articles))
	for i, article := range articles {
		responses[i] = *m.DomainToDTO(&article)
	}
	return responses
}

// DomainToShortDTO converts a domain Article entity to an ArticleShort DTO.
func (m *ArticleMapper) DomainToShortDTO(article *domain.Article) *dto.ArticleShort {
	return &dto.ArticleShort{
		ID:    article.ID,
		Title: article.Title,
		Date:  article.Date,
	}
}

// Domain to Model Conversions (Infrastructure layer)

// DomainToModel converts a domain Article entity to a persistence model.
func (m *ArticleMapper) DomainToModel(article *domain.Article) *model.Article {
	modelArticle := &model.Article{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		ImgBanner: article.ImgBanner,
		Date:      article.Date,
		SeasonID:  article.SeasonID,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}

	if article.Season != nil {
		seasonMapper := NewSeasonMapper()
		modelArticle.Season = seasonMapper.DomainToModel(article.Season)
	}

	return modelArticle
}

// ModelToDomain converts a persistence model to a domain Article entity.
func (m *ArticleMapper) ModelToDomain(model *model.Article) *domain.Article {
	domainArticle := &domain.Article{
		ID:        model.ID,
		Title:     model.Title,
		Content:   model.Content,
		ImgBanner: model.ImgBanner,
		Date:      model.Date,
		SeasonID:  model.SeasonID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	if model.Season != nil {
		seasonMapper := NewSeasonMapper()
		domainArticle.Season = seasonMapper.ModelToDomain(model.Season)
	}

	return domainArticle
}

// ModelListToDomain converts a slice of persistence models to domain Article entities.
func (m *ArticleMapper) ModelListToDomain(models []model.Article) []domain.Article {
	articles := make([]domain.Article, len(models))
	for i, model := range models {
		articles[i] = *m.ModelToDomain(&model)
	}
	return articles
}

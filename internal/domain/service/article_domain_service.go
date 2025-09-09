package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// ArticleDomainService contains the business logic for article operations.
// It operates on domain entities and implements business rules without external dependencies.
type ArticleDomainService struct {
	articleRepository domain.ArticleRepository
	seasonRepository  domain.SeasonRepository
}

func NewArticleDomainService(articleRepository domain.ArticleRepository, seasonRepository domain.SeasonRepository) *ArticleDomainService {
	return &ArticleDomainService{
		articleRepository: articleRepository,
		seasonRepository:  seasonRepository,
	}
}

// CreateArticle creates a new article after validating the referenced season exists.
func (s *ArticleDomainService) CreateArticle(ctx context.Context, article *domain.Article) error {
	// Verify that the season exists
	_, err := s.seasonRepository.GetSeasonByID(ctx, article.SeasonID)
	if err != nil {
		if err == constants.ErrRecordNotFound {
			return constants.ErrSeasonNotFound
		}
		return err
	}

	return s.articleRepository.CreateArticle(ctx, article)
}

// GetArticleByID retrieves an article by its ID.
func (s *ArticleDomainService) GetArticleByID(ctx context.Context, id uint64) (*domain.Article, error) {
	article, err := s.articleRepository.GetArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, constants.ErrRecordNotFound
	}
	return article, nil
}

// GetPaginatedArticles retrieves a paginated list of articles.
func (s *ArticleDomainService) GetPaginatedArticles(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Article, int64, error) {
	return s.articleRepository.GetPaginatedArticles(ctx, sort, order, page, pageSize)
}

// GetArticlesBySeasonID retrieves articles for a specific season after validating the season exists.
func (s *ArticleDomainService) GetArticlesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]domain.Article, int64, error) {
	// Verify that the season exists
	_, err := s.seasonRepository.GetSeasonByID(ctx, seasonID)
	if err != nil {
		if err == constants.ErrRecordNotFound {
			return nil, 0, constants.ErrSeasonNotFound
		}
		return nil, 0, err
	}

	return s.articleRepository.GetArticlesBySeasonID(ctx, seasonID, sort, order, page, pageSize)
}

// UpdateArticle updates an existing article after validating referenced entities exist.
func (s *ArticleDomainService) UpdateArticle(ctx context.Context, articleID uint64, updatedArticle *domain.Article) (*domain.Article, error) {
	// Check if article exists
	existingArticle, err := s.articleRepository.GetArticleByID(ctx, articleID)
	if err != nil {
		return nil, err
	}
	if existingArticle == nil {
		return nil, constants.ErrRecordNotFound
	}

	// If updating season, verify it exists
	if updatedArticle.SeasonID != 0 && updatedArticle.SeasonID != existingArticle.SeasonID {
		_, err := s.seasonRepository.GetSeasonByID(ctx, updatedArticle.SeasonID)
		if err != nil {
			if err == constants.ErrRecordNotFound {
				return nil, constants.ErrSeasonNotFound
			}
			return nil, err
		}
	}

	if err := s.articleRepository.UpdateArticle(ctx, articleID, updatedArticle); err != nil {
		return nil, err
	}

	// Return updated article
	return s.articleRepository.GetArticleByID(ctx, articleID)
}

// DeleteArticle deletes an article by its ID.
func (s *ArticleDomainService) DeleteArticle(ctx context.Context, id uint64) error {
	return s.articleRepository.DeleteArticle(ctx, id)
}

package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, article *model.Article) (*dto.ArticleShort, error)
	GetArticleByID(ctx context.Context, id uint64) (*model.Article, error)
	GetPaginatedArticles(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Article, int64, error)
	GetArticlesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]model.Article, int64, error)
	UpdateArticle(ctx context.Context, articleUpdate *dto.UpdateArticleRequest, articleID uint64) (*model.Article, error)
	DeleteArticle(ctx context.Context, id uint64) error
}

type articleService struct {
	ArticleRepository repository.ArticleRepository
	SeasonService     SeasonService
}

func NewArticleService(articleRepo repository.ArticleRepository, seasonService SeasonService) ArticleService {
	return &articleService{
		ArticleRepository: articleRepo,
		SeasonService:     seasonService,
	}
}

func (s *articleService) CreateArticle(ctx context.Context, article *model.Article) (*dto.ArticleShort, error) {
	// Verify that the season exists
	_, err := s.SeasonService.GetSeasonByID(ctx, article.SeasonID)
	if err != nil {
		if err == constants.ErrRecordNotFound {
			return nil, constants.ErrSeasonNotFound
		}
		return nil, fmt.Errorf("failed to verify season: %w", err)
	}

	if err := s.ArticleRepository.CreateArticle(ctx, article); err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return mapper.ToArticleShort(article), nil
}

func (s *articleService) GetArticleByID(ctx context.Context, id uint64) (*model.Article, error) {
	article, err := s.ArticleRepository.GetArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, constants.ErrRecordNotFound
	}
	return article, nil
}

func (s *articleService) GetPaginatedArticles(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Article, int64, error) {
	return s.ArticleRepository.GetPaginatedArticles(ctx, sort, order, page, pageSize)
}

func (s *articleService) GetArticlesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]model.Article, int64, error) {
	// Verify that the season exists
	_, err := s.SeasonService.GetSeasonByID(ctx, seasonID)
	if err != nil {
		if err == constants.ErrRecordNotFound {
			return nil, 0, constants.ErrSeasonNotFound
		}
		return nil, 0, fmt.Errorf("failed to verify season: %w", err)
	}

	return s.ArticleRepository.GetArticlesBySeasonID(ctx, seasonID, sort, order, page, pageSize)
}

func (s *articleService) UpdateArticle(ctx context.Context, articleUpdate *dto.UpdateArticleRequest, articleID uint64) (*model.Article, error) {
	article, err := s.ArticleRepository.GetArticleByID(ctx, articleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get article by ID: %w", err)
	}
	if article == nil {
		return nil, constants.ErrRecordNotFound
	}

	// If updating season, verify it exists
	if articleUpdate.SeasonID != nil {
		_, err := s.SeasonService.GetSeasonByID(ctx, *articleUpdate.SeasonID)
		if err != nil {
			if err == constants.ErrRecordNotFound {
				return nil, constants.ErrSeasonNotFound
			}
			return nil, fmt.Errorf("failed to verify season: %w", err)
		}
	}

	mapper.UpdateArticleFromDTO(article, articleUpdate)

	if err := s.ArticleRepository.UpdateArticle(ctx, articleID, article); err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return article, nil
}

func (s *articleService) DeleteArticle(ctx context.Context, id uint64) error {
	return s.ArticleRepository.DeleteArticle(ctx, id)
}

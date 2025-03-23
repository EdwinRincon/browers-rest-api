package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, article *model.Article) error
	GetArticleByID(ctx context.Context, id uint64) (*model.Article, error)
	GetAllArticles(ctx context.Context, page, pageSize uint64) ([]*model.Article, error)
	UpdateArticle(ctx context.Context, article *model.Article) error
	DeleteArticle(ctx context.Context, id uint64) error
}

type articleService struct {
	ArticleRepository repository.ArticleRepository
}

func NewArticleService(articleRepo repository.ArticleRepository) ArticleService {
	return &articleService{
		ArticleRepository: articleRepo,
	}
}

func (s *articleService) CreateArticle(ctx context.Context, article *model.Article) error {
	return s.ArticleRepository.CreateArticle(ctx, article)
}

func (s *articleService) GetArticleByID(ctx context.Context, id uint64) (*model.Article, error) {
	return s.ArticleRepository.GetArticleByID(ctx, id)
}

func (s *articleService) GetAllArticles(ctx context.Context, page, pageSize uint64) ([]*model.Article, error) {
	return s.ArticleRepository.GetAllArticles(ctx, page, pageSize)
}

func (s *articleService) UpdateArticle(ctx context.Context, article *model.Article) error {
	return s.ArticleRepository.UpdateArticle(ctx, article)
}

func (s *articleService) DeleteArticle(ctx context.Context, id uint64) error {
	return s.ArticleRepository.DeleteArticle(ctx, id)
}

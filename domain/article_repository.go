package domain

import (
	"context"
)

// ArticleRepository defines the interface for article persistence operations.
// This port belongs in the domain layer
type ArticleRepository interface {
	CreateArticle(ctx context.Context, article *Article) error
	GetArticleByID(ctx context.Context, id uint64) (*Article, error)
	GetPaginatedArticles(ctx context.Context, sort string, order string, page int, pageSize int) ([]Article, int64, error)
	GetArticlesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]Article, int64, error)
	UpdateArticle(ctx context.Context, id uint64, article *Article) error
	DeleteArticle(ctx context.Context, id uint64) error
}

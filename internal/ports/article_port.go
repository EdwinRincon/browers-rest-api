package ports

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
)

// ArticlePort defines the interface for article data access operations.
// This port separates the domain/application layer from the persistence adapter.
type ArticlePort interface {
	CreateArticle(ctx context.Context, article *model.Article) error
	GetArticleByID(ctx context.Context, id uint64) (*model.Article, error)
	GetPaginatedArticles(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Article, int64, error)
	GetArticlesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]model.Article, int64, error)
	UpdateArticle(ctx context.Context, id uint64, article *model.Article) error
	DeleteArticle(ctx context.Context, id uint64) error
}

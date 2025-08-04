package repository

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article *model.Article) error
	GetArticleByID(ctx context.Context, id uint64) (*model.Article, error)
	GetAllArticles(ctx context.Context, page, pageSize uint64) ([]*model.Article, error)
	UpdateArticle(ctx context.Context, article *model.Article) error
	DeleteArticle(ctx context.Context, id uint64) error
}

type ArticleRepositoryImpl struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &ArticleRepositoryImpl{db: db}
}

func (ar *ArticleRepositoryImpl) CreateArticle(ctx context.Context, article *model.Article) error {
	return ar.db.WithContext(ctx).Create(article).Error
}

func (ar *ArticleRepositoryImpl) GetArticleByID(ctx context.Context, id uint64) (*model.Article, error) {
	var article model.Article
	err := ar.db.WithContext(ctx).
		Preload("Season").
		First(&article, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrArticleNotFound
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}
	return &article, nil
}

func (ar *ArticleRepositoryImpl) GetAllArticles(ctx context.Context, page, pageSize uint64) ([]*model.Article, error) {
	if page < 1 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 25 // Default page size
	}
	if pageSize > 100 {
		pageSize = 100 // Maximum page size
	}

	offset := (page - 1) * pageSize
	// Guard against integer overflow
	if offset > math.MaxInt32 || pageSize > math.MaxInt32 {
		return nil, fmt.Errorf("pagination parameters too large")
	}

	var articles []*model.Article
	err := ar.db.WithContext(ctx).
		Order("date DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (ar *ArticleRepositoryImpl) UpdateArticle(ctx context.Context, article *model.Article) error {
	result := ar.db.WithContext(ctx).Save(article)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return constants.ErrArticleNotFound
	}
	return nil
}

func (ar *ArticleRepositoryImpl) DeleteArticle(ctx context.Context, id uint64) error {
	result := ar.db.WithContext(ctx).Delete(&model.Article{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return constants.ErrArticleNotFound
	}
	return nil
}

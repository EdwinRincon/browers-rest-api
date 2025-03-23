package repository

import (
	"context"
	"errors"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

var ErrArticleNotFound = errors.New("article not found")

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
	err := ar.db.WithContext(ctx).Preload("Season").First(&article, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return &article, nil
}

func (ar *ArticleRepositoryImpl) GetAllArticles(ctx context.Context, page, pageSize uint64) ([]*model.Article, error) {
	var articles []*model.Article
	offset := (page - 1) * pageSize
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
		return ErrArticleNotFound
	}
	return nil
}

func (ar *ArticleRepositoryImpl) DeleteArticle(ctx context.Context, id uint64) error {
	result := ar.db.WithContext(ctx).Delete(&model.Article{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrArticleNotFound
	}
	return nil
}

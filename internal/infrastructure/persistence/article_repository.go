package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/adapter/persistence"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

const (
	whereIDClause = "id = ?"
)

// ArticleRepositoryImpl implements domain.ArticleRepository interface.
type ArticleRepositoryImpl struct {
	db     *gorm.DB
	mapper *persistence.ArticlePersistenceMapper
}

// NewArticleRepository creates a new ArticleRepositoryImpl.
func NewArticleRepository(db *gorm.DB) *ArticleRepositoryImpl {
	return &ArticleRepositoryImpl{
		db:     db,
		mapper: persistence.NewArticlePersistenceMapper(),
	}
}

// GetArticleByID retrieves an article by its ID, preloading the Season.
func (ar *ArticleRepositoryImpl) GetArticleByID(ctx context.Context, id uint64) (*domain.Article, error) {
	var article model.Article
	result := ar.db.WithContext(ctx).
		Preload("Season").
		Where(whereIDClause, id).
		First(&article)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return ar.mapper.ModelToDomain(&article), nil
}

// GetPaginatedArticles retrieves a paginated list of articles with their seasons and total count.
func (ar *ArticleRepositoryImpl) GetPaginatedArticles(
	ctx context.Context,
	sort string,
	order string,
	page int,
	pageSize int,
) ([]domain.Article, int64, error) {
	var (
		articles []model.Article
		total    int64
	)

	// Count total records
	if err := ar.db.WithContext(ctx).Model(&model.Article{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total articles: %w", err)
	}

	// Build base query with eager loading
	query := ar.db.WithContext(ctx).
		Model(&model.Article{}).
		Preload("Season")

	// Apply sorting (safe and validated)
	col, raw, err := BuildOrderClause(EntityArticle, sort, order)
	if err != nil {
		return nil, 0, fmt.Errorf("error building sort clause: %w", err)
	}

	if raw != "" {
		query = query.Order(raw)
	} else {
		query = query.Order(col)
	}

	if pageSize <= 0 {
		pageSize = 10
	}
	if page < 0 {
		page = 0
	}
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute and map result
	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching paginated articles: %w", err)
	}

	return ar.mapper.ModelListToDomain(articles), total, nil
}

func (ar *ArticleRepositoryImpl) GetArticlesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]domain.Article, int64, error) {
	var articles []model.Article
	var total int64

	// Count total records for the season
	countQuery := ar.db.WithContext(ctx).Model(&model.Article{}).Where("season_id = ?", seasonID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting articles for season: %w", err)
	}

	query := ar.db.WithContext(ctx).Model(&model.Article{}).
		Preload("Season").
		Where("season_id = ?", seasonID)

	if sort != "" {
		col, raw, err := BuildOrderClause(EntityArticle, sort, order)
		if err != nil {
			return nil, 0, fmt.Errorf("error building sort clause: %w", err)
		}

		if raw != "" {
			query = query.Order(raw)
		} else {
			query = query.Order(col)
		}
	}

	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching articles for season: %w", err)
	}

	return ar.mapper.ModelListToDomain(articles), total, nil
}

func (ar *ArticleRepositoryImpl) CreateArticle(ctx context.Context, article *domain.Article) error {
	modelArticle := ar.mapper.DomainToModel(article)
	return ar.db.WithContext(ctx).Create(modelArticle).Error
}

func (ar *ArticleRepositoryImpl) UpdateArticle(ctx context.Context, id uint64, article *domain.Article) error {
	modelArticle := ar.mapper.DomainToModel(article)
	return ar.db.WithContext(ctx).
		Model(&model.Article{}).
		Where(whereIDClause, id).
		Select("*").
		Updates(modelArticle).Error
}

func (ar *ArticleRepositoryImpl) DeleteArticle(ctx context.Context, id uint64) error {
	return ar.db.WithContext(ctx).Delete(&model.Article{}, whereIDClause, id).Error
}

package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
)

func UpdateArticleFromDTO(article *model.Article, dto *dto.UpdateArticleRequest) {
	if dto.Title != nil {
		article.Title = *dto.Title
	}
	if dto.Content != nil {
		article.Content = *dto.Content
	}
	if dto.ImgBanner != nil {
		article.ImgBanner = *dto.ImgBanner
	}
	if dto.Date != nil {
		article.Date = *dto.Date
	}
	if dto.SeasonID != nil {
		article.SeasonID = *dto.SeasonID
	}
}

func ToArticleResponse(article *model.Article) *dto.ArticleResponse {
	var season dto.SeasonShort
	if article.Season != nil {
		seasonShort := ToSeasonShort(article.Season)
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

func ToArticleResponseList(articles []model.Article) []dto.ArticleResponse {
	articleResponses := make([]dto.ArticleResponse, len(articles))
	for i, article := range articles {
		articleResponses[i] = *ToArticleResponse(&article)
	}
	return articleResponses
}

func ToArticleShort(article *model.Article) *dto.ArticleShort {
	return &dto.ArticleShort{
		ID:    article.ID,
		Title: article.Title,
		Date:  article.Date,
	}
}

func ToArticle(dto *dto.CreateArticleRequest) *model.Article {
	return &model.Article{
		Title:     dto.Title,
		Content:   dto.Content,
		ImgBanner: dto.ImgBanner,
		Date:      dto.Date,
		SeasonID:  dto.SeasonID,
	}
}

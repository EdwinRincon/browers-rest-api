package dto

import (
	"time"
)

type CreateArticleRequest struct {
	Title     string    `json:"title" binding:"required,max=100"`
	Content   string    `json:"content" binding:"required"`
	ImgBanner string    `json:"img_banner,omitempty" binding:"omitempty,url"`
	Date      time.Time `json:"date" binding:"required"`
	SeasonID  uint64    `json:"season_id" binding:"required"`
}

type UpdateArticleRequest struct {
	Title     *string    `json:"title,omitempty" binding:"omitempty,max=100"`
	Content   *string    `json:"content,omitempty"`
	ImgBanner *string    `json:"img_banner,omitempty" binding:"omitempty,url"`
	Date      *time.Time `json:"date,omitempty"`
	SeasonID  *uint64    `json:"season_id,omitempty"`
}

type ArticleResponse struct {
	ID        uint64      `json:"id"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	ImgBanner string      `json:"img_banner,omitempty"`
	Date      time.Time   `json:"date"`
	Season    SeasonShort `json:"season,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type ArticleShort struct {
	ID    uint64    `json:"id"`
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
}

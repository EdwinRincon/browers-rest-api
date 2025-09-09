package domain

import (
	"time"
)

// Article represents the domain entity for an article.
// It encapsulates the business logic
type Article struct {
	ID        uint64
	Title     string
	Content   string
	ImgBanner string
	Date      time.Time
	SeasonID  uint64
	Season    *Season
	CreatedAt time.Time
	UpdatedAt time.Time
}

package model

import (
	"time"

	"gorm.io/gorm"
)

type Lineup struct {
	ID        uint64         `gorm:"primaryKey" json:"id" form:"id"`
	Position  string         `gorm:"type:varchar(5);not null" json:"position" form:"position" binding:"required,oneof=por ceni cend lati med latd del deli deld"`
	PlayerID  uint64         `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"player_id" form:"player_id" binding:"required"`
	Player    *Player        `gorm:"foreignKey:PlayerID" json:"player,omitempty" form:"player"`
	MatchID   uint64         `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"match_id" form:"match_id" binding:"required"`
	Match     *Match         `gorm:"foreignKey:MatchID" json:"match,omitempty" form:"match"`
	Starting  bool           `gorm:"default:false" json:"starting" form:"starting"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" form:"-" swaggerignore:"true"`
}

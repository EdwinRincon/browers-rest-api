package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Players struct {
	ID            uint64         `gorm:"primaryKey" json:"id" form:"id"`
	NickName      string         `gorm:"type:varchar(20)" json:"nick_name" form:"nick_name"`
	Height        float32        `json:"height" form:"height"`
	Country       string         `gorm:"type:varchar(10);not null" json:"country" form:"country"`
	Country2      string         `gorm:"type:varchar(10)" json:"country2" form:"country2"`
	Foot          string         `gorm:"type:varchar(1);not null" json:"foot" form:"foot"`
	Age           uint8          `gorm:"type:tinyint(3);not null" json:"age" form:"age"`
	SquadNumber   uint8          `gorm:"type:tinyint(3);not null" json:"squad_number" form:"squad_number"`
	Rating        uint8          `gorm:"type:tinyint(3);not null;default:0" json:"rating" form:"rating"`
	Matches       uint16         `gorm:"type:smallint(5);not null;default:0" json:"matches" form:"matches"`
	YCards        uint8          `gorm:"type:tinyint(3);not null;default:0" json:"y_cards" form:"y_cards"`
	RCards        uint8          `gorm:"type:tinyint(3);not null;default:0" json:"r_cards" form:"r_cards"`
	Goals         uint16         `gorm:"type:smallint(5);not null;default:0" json:"goals" form:"goals"`
	Assists       uint16         `gorm:"type:smallint(5);not null;default:0" json:"assists" form:"assists"`
	Saves         uint16         `gorm:"type:smallint(5);not null;default:0" json:"saves" form:"saves"`
	Position      string         `gorm:"type:varchar(5);not null" json:"position" form:"position"`
	Injured       bool           `gorm:"default:false" json:"injured" form:"injured"`
	CarrerSummary string         `gorm:"type:varchar(1000);not null" json:"carrer_summary" form:"carrer_summary"`
	Teams         []Teams        `gorm:"many2many:player_teams;" json:"teams" form:"teams"`
	MVPCount      uint8          `gorm:"type:tinyint(3);not null;default:0" json:"mvp_count" form:"mvp_count"`
	UsersID       string         `json:"users_id" form:"users_id"`
	Users         Users          `gorm:"foreignKey:UsersID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"users" form:"users"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}

// TODO: Pasar a BBDD como tabla nueva y hacer relación con Players
// ValidPositions contiene una lista de posiciones válidas
var validPositions = []string{"del", "mc", "def", "por"}

// PositionValidator es un validador personalizado para verificar si la posición es válida
var PositionValidator validator.Func = func(fl validator.FieldLevel) bool {
	position, ok := fl.Field().Interface().(string)
	if ok {
		for _, valid := range validPositions {
			if position == valid {
				return true
			}
		}
	}
	return false
}

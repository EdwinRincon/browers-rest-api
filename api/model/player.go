package model

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Players struct {
	gorm.Model
	ID            uint64 `gorm:"primaryKey"`
	NickName      string `gorm:"type:varchar(20)"`
	Height        float32
	Country       string `gorm:"type:varchar(10);not null"`
	Country2      string `gorm:"type:varchar(10)"`
	Foot          string `gorm:"type:varchar(1);not null"`
	Age           uint8  `gorm:"type:int(2);not null"`
	SquadNumber   uint8  `gorm:"type:int(2);not null"`
	Rating        uint8  `gorm:"type:int(2);not null;default:0"`
	Matches       uint16 `gorm:"type:int(3);not null;default:0"`
	YCards        uint8  `gorm:"type:int(2);not null;default:0"`
	RCards        uint8  `gorm:"type:int(2);not null;default:0"`
	Goals         uint16 `gorm:"type:int(3);not null;default:0"`
	Assists       uint16 `gorm:"type:int(3);not null;default:0"`
	Saves         uint16 `gorm:"type:int(3);not null;default:0"`
	Position      string `gorm:"type:varchar(5);not null"`
	Injured       bool   `gorm:"default:false"`
	CarrerSummary string `gorm:"type:varchar(1000);not null"`
	TeamsID       uint64 `form:"teams_id"`
	Teams         Teams
	Mvp           uint8  `gorm:"type:int(2);not null;default:0"`
	UsersID       uint64 `form:"users_id"`
	Users         Users
}

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

package domain

import "time"

// PlayerStat represents the statistics of a player in a specific match and season in the domain layer.
// This is a pure business entity
type PlayerStat struct {
	ID       uint64
	PlayerID uint64
	MatchID  uint64
	SeasonID uint64
	TeamID   *uint64

	Goals         uint8
	Assists       uint8
	Saves         uint8
	YellowCards   uint8
	RedCards      uint8
	Rating        uint8
	IsStarting    bool
	MinutesPlayed uint8
	IsMVP         bool
	Position      string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsValid performs basic domain validation for the player stat.
func (ps *PlayerStat) IsValid() bool {
	return ps.PlayerID > 0 &&
		ps.MatchID > 0 &&
		ps.SeasonID > 0 &&
		ps.Rating <= 100 &&
		ps.MinutesPlayed <= 120 && // max match duration + extra time
		ps.isValidPosition()
}

// isValidPosition checks if the position is one of the allowed values.
func (ps *PlayerStat) isValidPosition() bool {
	if ps.Position == "" {
		return true // Position can be empty
	}
	validPositions := map[string]bool{
		"por":  true, // portero
		"ceni": true, // central izquierdo
		"cend": true, // central derecho
		"lati": true, // lateral izquierdo
		"latd": true, // lateral derecho
		"med":  true, // mediocampista
		"del":  true, // delantero
		"deli": true, // delantero izquierdo
		"deld": true, // delantero derecho
	}
	return validPositions[ps.Position]
}

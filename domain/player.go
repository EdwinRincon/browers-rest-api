package domain

import "time"

// Player represents a football player in the domain layer.
// This is a pure business entity without infrastructure concerns.
type Player struct {
	ID               uint64
	NickName         string
	Height           uint16
	Country          string
	SecondaryCountry string
	Foot             string // L or R
	Age              uint8
	SquadNumber      uint8
	Rating           uint8
	Matches          uint16
	YCards           uint8 // Yellow cards
	RCards           uint8 // Red cards
	Goals            uint16
	Assists          uint16
	Saves            uint16
	Position         string // por ceni cenm cend lati med latd del deli deld
	Injured          bool
	CareerSummary    string
	MVPCount         uint8
	UserID           *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// IsValid performs basic domain validation for the player.
func (p *Player) IsValid() bool {
	return p.NickName != "" &&
		len(p.NickName) <= 20 &&
		p.Height >= 100 && p.Height <= 220 &&
		len(p.Country) == 2 &&
		(p.Foot == "L" || p.Foot == "R") &&
		p.Age >= 16 && p.Age <= 50 &&
		p.SquadNumber >= 1 && p.SquadNumber <= 99 &&
		p.Rating <= 100 &&
		p.isValidPosition() &&
		len(p.CareerSummary) <= 1000
}

// isValidPosition checks if the position is one of the allowed values.
func (p *Player) isValidPosition() bool {
	validPositions := map[string]bool{
		"por":  true, // portero
		"ceni": true, // central izquierdo
		"cenm": true, // central medio
		"cend": true, // central derecho
		"lati": true, // lateral izquierdo
		"latd": true, // lateral derecho
		"med":  true, // mediocampista
		"del":  true, // delantero
		"deli": true, // delantero izquierdo
		"deld": true, // delantero derecho
	}
	return validPositions[p.Position]
}

package domain

import "time"

// Lineup represents a lineup entry for a match in the domain layer.
// This entity contains only business-relevant fields
type Lineup struct {
	ID        uint64
	Position  string // por, ceni, cend, lati, med, latd, del, deli, deld
	PlayerID  uint64
	MatchID   uint64
	Starting  bool
	CreatedAt time.Time
	UpdatedAt time.Time

	// Related entities
	Player *Player
	Match  *Match
}

// IsValid performs basic domain validation for the lineup.
func (l *Lineup) IsValid() bool {
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

	return l.PlayerID > 0 &&
		l.MatchID > 0 &&
		validPositions[l.Position]
}

// IsStartingPlayer returns true if this lineup entry is for a starting player.
func (l *Lineup) IsStartingPlayer() bool {
	return l.Starting
}

// IsSubstitute returns true if this lineup entry is for a substitute player.
func (l *Lineup) IsSubstitute() bool {
	return !l.Starting
}

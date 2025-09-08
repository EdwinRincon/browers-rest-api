package domain

import (
	"context"
)

// PlayerRepository defines the interface for player persistence operations.
// This port belongs in the domain layer following hexagonal architecture.
type PlayerRepository interface {
	CreatePlayer(ctx context.Context, player *Player) error
	GetPlayerByID(ctx context.Context, id uint64) (*Player, error)
	GetPlayerByNickName(ctx context.Context, nickName string) (*Player, error)
	GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]Player, int64, error)
	UpdatePlayer(ctx context.Context, id uint64, player *Player) error
	DeletePlayer(ctx context.Context, id uint64) error
}

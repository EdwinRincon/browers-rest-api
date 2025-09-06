package ports

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
)

// PlayerPort defines the interface for player data access operations.
// This port separates the domain/application layer from the persistence adapter.
type PlayerPort interface {
	CreatePlayer(ctx context.Context, player *model.Player) error
	GetPlayerByID(ctx context.Context, id uint64) (*model.Player, error)
	GetPlayerByNickName(ctx context.Context, nickName string) (*model.Player, error)
	GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Player, int64, error)
	UpdatePlayer(ctx context.Context, id uint64, player *model.Player) error
	DeletePlayer(ctx context.Context, id uint64) error
}

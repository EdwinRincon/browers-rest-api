package ports

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence"
)

// PlayerTeamPort defines the interface for player-team association operations.
// This port separates the domain/application layer from the persistence adapter.
type PlayerTeamPort interface {
	Create(ctx context.Context, playerTeam *model.PlayerTeam) error
	GetByPlayerID(ctx context.Context, playerID uint64) ([]model.PlayerTeam, error)
	DeleteByPlayerID(ctx context.Context, playerID uint64) error
	GetPlayerTeamByID(ctx context.Context, id uint64) (*model.PlayerTeam, error)
	GetPlayerTeamsByTeamID(ctx context.Context, teamID uint64) ([]model.PlayerTeam, error)
	GetPlayerTeamsBySeasonID(ctx context.Context, seasonID uint64) ([]model.PlayerTeam, error)
	GetPaginatedPlayerTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.PlayerTeam, int64, error)
	UpdatePlayerTeam(ctx context.Context, playerTeam *model.PlayerTeam) error
	DeletePlayerTeam(ctx context.Context, id uint64) error
	CheckOverlappingDates(ctx context.Context, data persistence.OverlapCheckData) (bool, error)
}

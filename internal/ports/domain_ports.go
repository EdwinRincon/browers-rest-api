package ports

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/domain"
)

// PlayerPort defines the interface for player operations.
// This port works with pure domain entities instead of persistence models.
type PlayerPort interface {
	CreatePlayer(ctx context.Context, player *domain.Player) error
	GetPlayerByID(ctx context.Context, id uint64) (*domain.Player, error)
	GetPlayerByNickName(ctx context.Context, nickName string) (*domain.Player, error)
	GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Player, int64, error)
	UpdatePlayer(ctx context.Context, id uint64, player *domain.Player) error
	DeletePlayer(ctx context.Context, id uint64) error
}

// TeamPort defines the interface for team operations.
// This port works with pure domain entities instead of persistence models.
type TeamPort interface {
	CreateTeam(ctx context.Context, team *domain.Team) error
	GetTeamByID(ctx context.Context, id uint64) (*domain.Team, error)
	GetTeamByName(ctx context.Context, fullName string) (*domain.Team, error)
	GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Team, int64, error)
	UpdateTeam(ctx context.Context, team *domain.Team) error
	DeleteTeam(ctx context.Context, id uint64) error
}


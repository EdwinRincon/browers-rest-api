package ports

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/domain"
)

// PlayerDomainPort defines the interface for player domain operations.
// This port works with pure domain entities instead of persistence models.
type PlayerDomainPort interface {
	CreatePlayer(ctx context.Context, player *domain.Player) error
	GetPlayerByID(ctx context.Context, id uint64) (*domain.Player, error)
	GetPlayerByNickName(ctx context.Context, nickName string) (*domain.Player, error)
	GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Player, int64, error)
	UpdatePlayer(ctx context.Context, id uint64, player *domain.Player) error
	DeletePlayer(ctx context.Context, id uint64) error
}

// TeamDomainPort defines the interface for team domain operations.
// This port works with pure domain entities instead of persistence models.
type TeamDomainPort interface {
	CreateTeam(ctx context.Context, team *domain.Team) error
	GetTeamByID(ctx context.Context, id uint64) (*domain.Team, error)
	GetTeamByName(ctx context.Context, fullName string) (*domain.Team, error)
	GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Team, int64, error)
	UpdateTeam(ctx context.Context, team *domain.Team) error
	DeleteTeam(ctx context.Context, id uint64) error
}

// SeasonDomainPort defines the interface for season domain operations.
// This port works with pure domain entities instead of persistence models.
type SeasonDomainPort interface {
	CreateSeason(ctx context.Context, season *domain.Season) error
	GetSeasonByID(ctx context.Context, id uint64) (*domain.Season, error)
	GetSeasonByYear(ctx context.Context, year uint16) (*domain.Season, error)
	GetCurrentSeason(ctx context.Context) (*domain.Season, error)
	GetPaginatedSeasons(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Season, int64, error)
	UpdateSeason(ctx context.Context, id uint64, season *domain.Season) error
	DeleteSeason(ctx context.Context, id uint64) error
	SetCurrentSeason(ctx context.Context, id uint64) error
}

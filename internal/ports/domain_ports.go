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

// SeasonPort defines the interface for season operations.
// This port works with pure domain entities instead of persistence models.
type SeasonPort interface {
	CreateSeason(ctx context.Context, season *domain.Season) error
	GetSeasonByID(ctx context.Context, id uint64) (*domain.Season, error)
	GetSeasonByYear(ctx context.Context, year uint16) (*domain.Season, error)
	GetCurrentSeason(ctx context.Context) (*domain.Season, error)
	GetPaginatedSeasons(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Season, int64, error)
	UpdateSeason(ctx context.Context, id uint64, season *domain.Season) error
	DeleteSeason(ctx context.Context, id uint64) error
	SetCurrentSeason(ctx context.Context, id uint64) error
}

// UserPort defines the interface for user operations.
// This port works with pure domain entities and serves both:
// - Domain services (inbound port for handlers)
// - Persistence adapters (outbound port for data access)
type UserPort interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.User, int64, error)
	UpdateUser(ctx context.Context, id string, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

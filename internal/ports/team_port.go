package ports

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
)

// TeamPort defines the interface for team data access operations.
// This port separates the domain/application layer from the persistence adapter.
type TeamPort interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeamByID(ctx context.Context, id uint64) (*model.Team, error)
	GetTeamByName(ctx context.Context, fullName string) (*model.Team, error)
	GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Team, int64, error)
	UpdateTeam(ctx context.Context, team *model.Team) error
	DeleteTeam(ctx context.Context, id uint64) error
}

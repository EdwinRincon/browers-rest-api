package domain

import (
	"context"
)

// TeamRepository defines the interface for team persistence operations.
// This interface belongs in the domain layer following hexagonal architecture.
// All methods preserved from api/repository/team_repository.go interface.
type TeamRepository interface {
	CreateTeam(ctx context.Context, team *Team) error
	GetTeamByID(ctx context.Context, id uint64) (*Team, error)
	GetTeamByName(ctx context.Context, fullName string) (*Team, error)
	GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]Team, int64, error)
	UpdateTeam(ctx context.Context, team *Team) error
	DeleteTeam(ctx context.Context, id uint64) error
}

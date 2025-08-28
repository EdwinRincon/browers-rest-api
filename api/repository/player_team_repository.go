package repository

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

// PlayerTeamRepository defines the interface for player-team association operations
type PlayerTeamRepository interface {
	Create(ctx context.Context, playerTeam *model.PlayerTeam) error
	DeleteByPlayerAndTeam(ctx context.Context, playerID uint64, teamID uint64) error
	GetByPlayerID(ctx context.Context, playerID uint64) ([]model.PlayerTeam, error)
	DeleteByPlayerID(ctx context.Context, playerID uint64) error
}

// PlayerTeamRepositoryImpl implements PlayerTeamRepository
type PlayerTeamRepositoryImpl struct {
	db *gorm.DB
}

// NewPlayerTeamRepository creates a new PlayerTeamRepository
func NewPlayerTeamRepository(db *gorm.DB) PlayerTeamRepository {
	return &PlayerTeamRepositoryImpl{
		db: db,
	}
}

// Create adds a new player-team association
func (ptr *PlayerTeamRepositoryImpl) Create(ctx context.Context, playerTeam *model.PlayerTeam) error {
	if err := ptr.db.WithContext(ctx).Create(playerTeam).Error; err != nil {
		return fmt.Errorf("failed to create player team association: %w", err)
	}
	return nil
}

// DeleteByPlayerAndTeam removes a specific player-team association
func (ptr *PlayerTeamRepositoryImpl) DeleteByPlayerAndTeam(ctx context.Context, playerID uint64, teamID uint64) error {
	if err := ptr.db.WithContext(ctx).
		Where("player_id = ? AND team_id = ?", playerID, teamID).
		Delete(&model.PlayerTeam{}).Error; err != nil {
		return fmt.Errorf("failed to delete player team association: %w", err)
	}
	return nil
}

// GetByPlayerID gets all team associations for a player
func (ptr *PlayerTeamRepositoryImpl) GetByPlayerID(ctx context.Context, playerID uint64) ([]model.PlayerTeam, error) {
	var playerTeams []model.PlayerTeam
	if err := ptr.db.WithContext(ctx).
		Preload("Team").
		Where("player_id = ?", playerID).
		Find(&playerTeams).Error; err != nil {
		return nil, fmt.Errorf("failed to get player team associations: %w", err)
	}
	return playerTeams, nil
}

// DeleteByPlayerID removes all team associations for a player
func (ptr *PlayerTeamRepositoryImpl) DeleteByPlayerID(ctx context.Context, playerID uint64) error {
	if err := ptr.db.WithContext(ctx).
		Where("player_id = ?", playerID).
		Delete(&model.PlayerTeam{}).Error; err != nil {
		return fmt.Errorf("failed to delete player team associations: %w", err)
	}
	return nil
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

const (
	// PreloadPlayerTeamsTeam is the relation path to preload the Team via PlayerTeam
	PreloadPlayerTeamsTeam = "PlayerTeams.Team"

	PreloadUser         = "User"
	WhereIDEquals       = "id = ?"
	WhereNickNameEquals = "nick_name = ?"
)

type PlayerRepository interface {
	CreatePlayer(ctx context.Context, player *model.Player) error
	GetActivePlayerByID(ctx context.Context, id uint64) (*model.Player, error)
	GetUnscopedPlayerByID(ctx context.Context, id uint64) (*model.Player, error)
	GetActivePlayerByNickName(ctx context.Context, nickName string) (*model.Player, error)
	GetUnscopedPlayerByNickName(ctx context.Context, nickName string) (*model.Player, error)
	GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Player, int64, error)
	UpdatePlayer(ctx context.Context, id uint64, player *model.Player) error
	DeletePlayer(ctx context.Context, id uint64) error
}

type PlayerRepositoryImpl struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &PlayerRepositoryImpl{db: db}
}

// GetActivePlayerByNickName retrieves an active (not deleted) player by their nickname.
func (pr *PlayerRepositoryImpl) GetActivePlayerByNickName(ctx context.Context, nickName string) (*model.Player, error) {
	var player model.Player
	result := pr.db.WithContext(ctx).
		Preload(PreloadUser).
		Preload(PreloadPlayerTeamsTeam).
		Where(WhereNickNameEquals, nickName).
		First(&player)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting player by nickname: %w", result.Error)
	}
	return &player, nil
}

// GetUnscopedPlayerByNickName retrieves a player by their nickname, including soft-deleted records.
func (pr *PlayerRepositoryImpl) GetUnscopedPlayerByNickName(ctx context.Context, nickName string) (*model.Player, error) {
	var player model.Player
	result := pr.db.WithContext(ctx).
		Preload(PreloadUser).
		Preload(PreloadPlayerTeamsTeam).
		Unscoped().
		Where(WhereNickNameEquals, nickName).
		First(&player)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &player, result.Error
}

// GetActivePlayerByID retrieves an active (not deleted) player by their ID, with preloaded relations.
func (pr *PlayerRepositoryImpl) GetActivePlayerByID(ctx context.Context, id uint64) (*model.Player, error) {
	var player model.Player
	result := pr.db.WithContext(ctx).
		Preload(PreloadUser).
		Preload(PreloadPlayerTeamsTeam).
		Where(WhereIDEquals, id).
		First(&player)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &player, result.Error
}

// GetUnscopedPlayerByID retrieves a player by their ID including soft-deleted records, with preloaded relations.
func (pr *PlayerRepositoryImpl) GetUnscopedPlayerByID(ctx context.Context, id uint64) (*model.Player, error) {
	var player model.Player
	result := pr.db.WithContext(ctx).
		Preload(PreloadUser).
		Preload(PreloadPlayerTeamsTeam).
		Unscoped().
		Where(WhereIDEquals, id).
		First(&player)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &player, result.Error
}

// GetPaginatedPlayers retrieves a paginated list of players with their teams, user and total count.
func (pr *PlayerRepositoryImpl) GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Player, int64, error) {
	var players []model.Player
	var total int64

	// Count total records
	countQuery := pr.db.WithContext(ctx).Model(&model.Player{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total players: %w", err)
	}

	// Build the data query with eager loading
	query := pr.db.WithContext(ctx).Model(&model.Player{}).
		Preload(PreloadUser).
		Preload(PreloadPlayerTeamsTeam)

	// Apply sorting if provided
	if sort != "" && (order == "asc" || order == "desc") {
		query = query.Order(fmt.Sprintf("%s %s", sort, order))
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&players).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching players: %w", err)
	}

	return players, total, nil
}

func (pr *PlayerRepositoryImpl) CreatePlayer(ctx context.Context, player *model.Player) error {
	return pr.db.WithContext(ctx).Create(player).Error
}

func (pr *PlayerRepositoryImpl) UpdatePlayer(ctx context.Context, id uint64, player *model.Player) error {
	return pr.db.WithContext(ctx).
		Model(&model.Player{}).
		Where(WhereIDEquals, id).
		Select("*").
		Updates(player).Error
}

func (pr *PlayerRepositoryImpl) DeletePlayer(ctx context.Context, id uint64) error {
	return pr.db.WithContext(ctx).Delete(&model.Player{}, id).Error
}

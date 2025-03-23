package repository

import (
	"context"
	"errors"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

var ErrPlayerNotFound = errors.New("player not found")

type PlayerRepository interface {
	CreatePlayer(ctx context.Context, player *model.Player) error
	GetPlayerByID(ctx context.Context, id uint64) (*model.Player, error)
	GetAllPlayers(ctx context.Context, page uint64) ([]*model.Player, error)
	UpdatePlayer(ctx context.Context, player *model.Player) error
	DeletePlayer(ctx context.Context, id uint64) error
}

type PlayerRepositoryImpl struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &PlayerRepositoryImpl{db: db}
}

func (pr *PlayerRepositoryImpl) CreatePlayer(ctx context.Context, player *model.Player) error {
	return pr.db.WithContext(ctx).Create(player).Error
}

func (pr *PlayerRepositoryImpl) GetPlayerByID(ctx context.Context, id uint64) (*model.Player, error) {
	var player model.Player
	err := pr.db.WithContext(ctx).First(&player, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}
	return &player, nil
}

func (pr *PlayerRepositoryImpl) GetAllPlayers(ctx context.Context, page uint64) ([]*model.Player, error) {
	var players []*model.Player
	offset := (page - 1) * 10
	err := pr.db.WithContext(ctx).Offset(int(offset)).Limit(10).Find(&players).Error
	if err != nil {
		return nil, err
	}
	return players, nil
}

func (pr *PlayerRepositoryImpl) UpdatePlayer(ctx context.Context, player *model.Player) error {
	result := pr.db.WithContext(ctx).Save(player)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrPlayerNotFound
	}
	return nil
}

func (pr *PlayerRepositoryImpl) DeletePlayer(ctx context.Context, id uint64) error {
	result := pr.db.WithContext(ctx).Delete(&model.Player{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrPlayerNotFound
	}
	return nil
}

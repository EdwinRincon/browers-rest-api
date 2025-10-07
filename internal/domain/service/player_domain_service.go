package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// PlayerDomainService implements business logic for Player operations.
// It contains domain rules and validation while being infrastructure-agnostic.
type PlayerDomainService struct {
	playerRepository domain.PlayerRepository
}

func NewPlayerDomainService(playerRepository domain.PlayerRepository) *PlayerDomainService {
	return &PlayerDomainService{
		playerRepository: playerRepository,
	}
}

func (s *PlayerDomainService) CreatePlayer(ctx context.Context, player *domain.Player) error {
	// Validate domain rules
	if !player.IsValid() {
		return constants.ErrInvalidData
	}

	// Business rule: Check if player with same nickname already exists
	existingPlayer, err := s.playerRepository.GetPlayerByNickName(ctx, player.NickName)
	if err != nil {
		return fmt.Errorf("failed to check existing player: %w", err)
	}
	if existingPlayer != nil {
		return constants.ErrRecordAlreadyExists
	}

	return s.playerRepository.CreatePlayer(ctx, player)
}

func (s *PlayerDomainService) GetPlayerByID(ctx context.Context, id uint64) (*domain.Player, error) {
	if id == 0 {
		return nil, constants.ErrInvalidData
	}

	player, err := s.playerRepository.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player by ID: %w", err)
	}

	if player == nil {
		return nil, constants.ErrRecordNotFound
	}

	return player, nil
}

func (s *PlayerDomainService) GetPlayerByNickName(ctx context.Context, nickName string) (*domain.Player, error) {
	if nickName == "" {
		return nil, constants.ErrInvalidData
	}

	player, err := s.playerRepository.GetPlayerByNickName(ctx, nickName)
	if err != nil {
		return nil, fmt.Errorf("failed to get player by nickname: %w", err)
	}

	if player == nil {
		return nil, constants.ErrRecordNotFound
	}

	return player, nil
}

// GetPaginatedPlayers retrieves a paginated list of players.
func (s *PlayerDomainService) GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Player, int64, error) {
	return s.playerRepository.GetPaginatedPlayers(ctx, sort, order, page, pageSize)
}

// UpdatePlayer updates an existing player with business rule validation.
func (s *PlayerDomainService) UpdatePlayer(ctx context.Context, id uint64, player *domain.Player) (*domain.Player, error) {
	// Get existing player
	existingPlayer, err := s.playerRepository.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player by ID: %w", err)
	}
	if existingPlayer == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Apply updates to existing player
	if player.NickName != "" && player.NickName != existingPlayer.NickName {
		// Check if new nickname already exists
		nickNameCheck, err := s.playerRepository.GetPlayerByNickName(ctx, player.NickName)
		if err != nil {
			return nil, fmt.Errorf("failed to check nickname availability: %w", err)
		}
		if nickNameCheck != nil {
			return nil, constants.ErrRecordAlreadyExists
		}
		existingPlayer.NickName = player.NickName
	}
	if player.Position != "" {
		existingPlayer.Position = player.Position
	}
	if player.Height != 0 {
		existingPlayer.Height = player.Height
	}
	if player.Country != "" {
		existingPlayer.Country = player.Country
	}
	if player.SecondaryCountry != "" {
		existingPlayer.SecondaryCountry = player.SecondaryCountry
	}
	if player.Foot != "" {
		existingPlayer.Foot = player.Foot
	}
	if player.Age != 0 {
		existingPlayer.Age = player.Age
	}
	if player.SquadNumber != 0 {
		existingPlayer.SquadNumber = player.SquadNumber
	}
	if player.Rating != 0 {
		existingPlayer.Rating = player.Rating
	}
	if player.CareerSummary != "" {
		existingPlayer.CareerSummary = player.CareerSummary
	}

	// Validate updated player
	if !existingPlayer.IsValid() {
		return nil, constants.ErrInvalidData
	}

	// Update the player
	err = s.playerRepository.UpdatePlayer(ctx, id, existingPlayer)
	if err != nil {
		return nil, fmt.Errorf("failed to update player: %w", err)
	}

	return existingPlayer, nil
}

func (s *PlayerDomainService) DeletePlayer(ctx context.Context, id uint64) error {
	if id == 0 {
		return constants.ErrInvalidData
	}

	// Check if player exists
	existingPlayer, err := s.playerRepository.GetPlayerByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get player by ID: %w", err)
	}
	if existingPlayer == nil {
		return constants.ErrRecordNotFound
	}

	return s.playerRepository.DeletePlayer(ctx, id)
}

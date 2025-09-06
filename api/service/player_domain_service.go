package service

import (
	"context"
	"fmt"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

// PlayerDomainService represents the domain-focused player service.
// This service works with domain entities and demonstrates hexagonal architecture principles.
type PlayerDomainService interface {
	CreatePlayer(ctx context.Context, createRequest *dto.CreatePlayerRequest) (*dto.PlayerShort, error)
	GetPlayerByID(ctx context.Context, id uint64) (*domain.Player, error)
	GetPlayerByNickName(ctx context.Context, nickName string) (*domain.Player, error)
	GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Player, int64, error)
	UpdatePlayer(ctx context.Context, playerUpdate *dto.UpdatePlayerRequest, playerID uint64) (*domain.Player, error)
	DeletePlayer(ctx context.Context, id uint64) error

	// Domain-specific business methods
	GetEligiblePlayers(ctx context.Context) ([]domain.Player, error)
	GetPlayersByPosition(ctx context.Context, position string) ([]domain.Player, error)
}

type playerDomainService struct {
	playerPort     ports.PlayerDomainPort
	seasonPort     ports.SeasonDomainPort
	playerTeamPort ports.PlayerTeamPort // Still using persistence model for now
}

// NewPlayerDomainService creates a new domain-focused player service.
func NewPlayerDomainService(
	playerPort ports.PlayerDomainPort,
	seasonPort ports.SeasonDomainPort,
	playerTeamPort ports.PlayerTeamPort,
) PlayerDomainService {
	return &playerDomainService{
		playerPort:     playerPort,
		seasonPort:     seasonPort,
		playerTeamPort: playerTeamPort,
	}
}

// getCurrentSeason fetches the current active season using domain entities
func (s *playerDomainService) getCurrentSeason(ctx context.Context) (*domain.Season, error) {
	season, err := s.seasonPort.GetCurrentSeason(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current season: %w", err)
	}

	if season == nil {
		return nil, fmt.Errorf("no current season is set in the system")
	}

	return season, nil
}

// validatePlayerBusinessRules applies domain-specific business rules
func (s *playerDomainService) validatePlayerBusinessRules(player *domain.Player) error {
	// Check basic domain validation first
	if !player.IsValid() {
		// Check specific validation failures for better error messages
		if player.SquadNumber < 1 || player.SquadNumber > 99 {
			return fmt.Errorf("squad number must be between 1 and 99")
		}
		if player.Foot != "L" && player.Foot != "R" {
			return fmt.Errorf("foot must be either L or R")
		}
		if player.NickName == "" {
			return fmt.Errorf("nickname is required")
		}
		if len(player.NickName) > 20 {
			return fmt.Errorf("nickname must be 20 characters or less")
		}
		return fmt.Errorf("player fails domain validation")
	}

	// Example business rule: Goalkeepers must have 0 goals and assists
	if player.Position == "por" && (player.Goals > 0 || player.Assists > 0) {
		return fmt.Errorf("goalkeepers cannot have goals or assists")
	}

	return nil
}

// checkNicknameExistsInDomain checks if a player with the given nickname already exists using domain entities
func (s *playerDomainService) checkNicknameExistsInDomain(ctx context.Context, nickname string) error {
	player, err := s.playerPort.GetPlayerByNickName(ctx, nickname)
	if err != nil {
		return fmt.Errorf("failed to check existing player: %w", err)
	}
	if player != nil {
		return constants.ErrRecordAlreadyExists
	}
	return nil
}

func (s *playerDomainService) CreatePlayer(ctx context.Context, createRequest *dto.CreatePlayerRequest) (*dto.PlayerShort, error) {
	// Convert DTO to domain entity
	player := mapper.CreatePlayerRequestToDomain(createRequest)

	// Apply domain validation
	if err := s.validatePlayerBusinessRules(player); err != nil {
		return nil, err
	}

	// Check business rule: nickname uniqueness
	if err := s.checkNicknameExistsInDomain(ctx, player.NickName); err != nil {
		return nil, err
	}

	// Set creation timestamp
	player.CreatedAt = time.Now()
	player.UpdatedAt = time.Now()

	// Persist via domain port
	if err := s.playerPort.CreatePlayer(ctx, player); err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	// Convert domain entity back to DTO for response
	return mapper.PlayerDomainToShort(player), nil
}

func (s *playerDomainService) GetPlayerByNickName(ctx context.Context, nickName string) (*domain.Player, error) {
	player, err := s.playerPort.GetPlayerByNickName(ctx, nickName)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, constants.ErrRecordNotFound
	}
	return player, nil
}

func (s *playerDomainService) GetPlayerByID(ctx context.Context, id uint64) (*domain.Player, error) {
	player, err := s.playerPort.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, constants.ErrRecordNotFound
	}
	return player, nil
}

func (s *playerDomainService) GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Player, int64, error) {
	return s.playerPort.GetPaginatedPlayers(ctx, sort, order, page, pageSize)
}

func (s *playerDomainService) UpdatePlayer(ctx context.Context, playerUpdate *dto.UpdatePlayerRequest, playerID uint64) (*domain.Player, error) {
	// Get existing player from domain
	player, err := s.GetPlayerByID(ctx, playerID)
	if err != nil {
		return nil, err
	}

	// Check nickname uniqueness if it's being changed
	if playerUpdate.NickName != nil && *playerUpdate.NickName != player.NickName {
		if err := s.checkNicknameExistsInDomain(ctx, *playerUpdate.NickName); err != nil {
			return nil, err
		}
	}

	// Apply updates to domain entity
	mapper.UpdatePlayerRequestToDomain(player, playerUpdate)

	// Apply domain validation
	if err := s.validatePlayerBusinessRules(player); err != nil {
		return nil, err
	}

	// Update timestamp
	player.UpdatedAt = time.Now()

	// Persist via domain port
	if err := s.playerPort.UpdatePlayer(ctx, playerID, player); err != nil {
		return nil, fmt.Errorf("failed to update player: %w", err)
	}

	return player, nil
}

func (s *playerDomainService) DeletePlayer(ctx context.Context, id uint64) error {
	// Check if player exists (domain validation)
	_, err := s.GetPlayerByID(ctx, id)
	if err != nil {
		return err
	}

	// Additional business rules could be applied here
	// For example: check if player is in active lineups, has stats, etc.

	return s.playerPort.DeletePlayer(ctx, id)
}

// Additional domain methods that showcase business logic

// GetEligiblePlayers returns players who are eligible to play (not injured, valid squad number)
func (s *playerDomainService) GetEligiblePlayers(ctx context.Context) ([]domain.Player, error) {
	allPlayers, _, err := s.playerPort.GetPaginatedPlayers(ctx, "", "", 0, 1000) // Get all players
	if err != nil {
		return nil, err
	}

	var eligiblePlayers []domain.Player
	for _, player := range allPlayers {
		// Check if player is eligible: not injured and has valid squad number
		if !player.Injured && player.SquadNumber > 0 && player.SquadNumber <= 99 {
			eligiblePlayers = append(eligiblePlayers, player)
		}
	}

	return eligiblePlayers, nil
}

// GetPlayersByPosition returns players filtered by position using domain logic
func (s *playerDomainService) GetPlayersByPosition(ctx context.Context, position string) ([]domain.Player, error) {
	allPlayers, _, err := s.playerPort.GetPaginatedPlayers(ctx, "", "", 0, 1000)
	if err != nil {
		return nil, err
	}

	var filteredPlayers []domain.Player
	for _, player := range allPlayers {
		if player.Position == position {
			filteredPlayers = append(filteredPlayers, player)
		}
	}

	return filteredPlayers, nil
}

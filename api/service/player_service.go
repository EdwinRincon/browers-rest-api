package service

import (
	"context"
	"fmt"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type PlayerService interface {
	CreatePlayer(ctx context.Context, createRequest *dto.CreatePlayerRequest) (*dto.PlayerShort, error)
	GetPlayerByID(ctx context.Context, id uint64) (*model.Player, error)
	GetPlayerByNickName(ctx context.Context, nickName string) (*model.Player, error)
	GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Player, int64, error)
	UpdatePlayer(ctx context.Context, playerUpdate *dto.UpdatePlayerRequest, playerID uint64) (*model.Player, error)
	DeletePlayer(ctx context.Context, id uint64) error
}

type playerService struct {
	PlayerRepository     repository.PlayerRepository
	PlayerTeamRepository repository.PlayerTeamRepository
	SeasonRepository     repository.SeasonRepository
}

func NewPlayerService(playerRepo repository.PlayerRepository, playerTeamRepo repository.PlayerTeamRepository, seasonRepo repository.SeasonRepository) PlayerService {
	return &playerService{
		PlayerRepository:     playerRepo,
		PlayerTeamRepository: playerTeamRepo,
		SeasonRepository:     seasonRepo,
	}
}

// getCurrentSeasonID fetches the current active season ID
func (s *playerService) getCurrentSeasonID(ctx context.Context) (uint64, error) {
	// Get the current season from the repository
	season, err := s.SeasonRepository.GetCurrentSeason(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get current season: %w", err)
	}

	// If no current season is set, return an error
	if season == nil {
		return 0, fmt.Errorf("no current season is set in the system")
	}

	return season.ID, nil
}

// addTeamAssociations adds team associations for a player
func (s *playerService) addTeamAssociations(ctx context.Context, playerID uint64, teamIDs []uint64) error {
	// Get the current season ID
	seasonID, err := s.getCurrentSeasonID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current season ID: %w", err)
	}

	for _, teamID := range teamIDs {
		playerTeam := &model.PlayerTeam{
			PlayerID:  playerID,
			TeamID:    teamID,
			SeasonID:  seasonID,
			StartDate: time.Now(),
		}
		if err := s.PlayerTeamRepository.Create(ctx, playerTeam); err != nil {
			return fmt.Errorf("failed to add team association: %w", err)
		}
	}
	return nil
}

// restorePlayer updates a soft-deleted player and restores it
func (s *playerService) restorePlayer(ctx context.Context, existing *model.Player, player *model.Player) error {
	// Update all relevant fields from the new player
	existing.NickName = player.NickName
	existing.Height = player.Height
	existing.Country = player.Country
	existing.Country2 = player.Country2
	existing.Foot = player.Foot
	existing.Age = player.Age
	existing.SquadNumber = player.SquadNumber
	existing.Position = player.Position
	existing.CareerSummary = player.CareerSummary
	existing.UserID = player.UserID

	// Reset DeletedAt to restore the player
	existing.DeletedAt.Valid = false
	existing.DeletedAt.Time = time.Time{}

	// Update the player with all fields
	return s.PlayerRepository.UpdatePlayer(ctx, existing.ID, existing)
}

// checkNicknameExists checks if a player with the given nickname already exists
func (s *playerService) checkNicknameExists(ctx context.Context, nickname string) error {
	activePlayer, err := s.PlayerRepository.GetActivePlayerByNickName(ctx, nickname)
	if err != nil {
		return fmt.Errorf("failed to check existing active player: %w", err)
	}
	if activePlayer != nil {
		return constants.ErrRecordAlreadyExists
	}
	return nil
}

// createNewPlayer creates a new player in the database
func (s *playerService) createNewPlayer(ctx context.Context, player *model.Player, teamIDs []uint64) (*dto.PlayerShort, error) {
	if err := s.PlayerRepository.CreatePlayer(ctx, player); err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	// If team IDs are provided, add team associations
	if len(teamIDs) > 0 {
		if err := s.addTeamAssociations(ctx, player.ID, teamIDs); err != nil {
			return nil, err
		}
	}

	return mapper.ToPlayerShort(player), nil
}

func (s *playerService) CreatePlayer(ctx context.Context, createRequest *dto.CreatePlayerRequest) (*dto.PlayerShort, error) {
	// Check if nickname already exists
	if err := s.checkNicknameExists(ctx, createRequest.NickName); err != nil {
		return nil, err
	}

	// Convert DTO to model
	player := mapper.ToPlayer(createRequest)

	// Check for soft-deleted player with the same nickname
	existing, err := s.PlayerRepository.GetUnscopedPlayerByNickName(ctx, createRequest.NickName)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing player: %w", err)
	}

	// If soft-deleted player exists, restore it
	if existing != nil && existing.DeletedAt.Valid {
		if err := s.restorePlayer(ctx, existing, player); err != nil {
			return nil, fmt.Errorf("failed to restore and update player: %w", err)
		}

		if len(createRequest.TeamIDs) > 0 {
			if err := s.addTeamAssociations(ctx, existing.ID, createRequest.TeamIDs); err != nil {
				return nil, err
			}
		}

		return mapper.ToPlayerShort(existing), nil
	}

	// Create new player if no soft-deleted player exists
	return s.createNewPlayer(ctx, player, createRequest.TeamIDs)
}

func (s *playerService) GetPlayerByNickName(ctx context.Context, nickName string) (*model.Player, error) {
	player, err := s.PlayerRepository.GetActivePlayerByNickName(ctx, nickName)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, constants.ErrRecordNotFound
	}
	return player, nil
}

func (s *playerService) GetPlayerByID(ctx context.Context, id uint64) (*model.Player, error) {
	player, err := s.PlayerRepository.GetActivePlayerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, constants.ErrRecordNotFound
	}
	return player, nil
}

func (s *playerService) GetPaginatedPlayers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Player, int64, error) {
	return s.PlayerRepository.GetPaginatedPlayers(ctx, sort, order, page, pageSize)
}

// updatePlayerTeams updates the team associations for a player
func (s *playerService) updatePlayerTeams(ctx context.Context, playerID uint64, teamIDs []uint64) error {
	// First, delete existing team associations
	if err := s.PlayerTeamRepository.DeleteByPlayerID(ctx, playerID); err != nil {
		return fmt.Errorf("failed to remove existing team associations: %w", err)
	}

	// Then create new team associations
	return s.addTeamAssociations(ctx, playerID, teamIDs)
}

func (s *playerService) UpdatePlayer(ctx context.Context, playerUpdate *dto.UpdatePlayerRequest, playerID uint64) (*model.Player, error) {
	player, err := s.PlayerRepository.GetActivePlayerByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player by ID: %w", err)
	}
	if player == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Check for nickname uniqueness if it's being changed
	if playerUpdate.NickName != nil && *playerUpdate.NickName != player.NickName {
		dup, err := s.PlayerRepository.GetActivePlayerByNickName(ctx, *playerUpdate.NickName)
		if err != nil {
			return nil, fmt.Errorf("failed to check duplicate nickname: %w", err)
		}
		if dup != nil && dup.ID != player.ID {
			return nil, constants.ErrRecordAlreadyExists
		}
	}

	// Update player data from DTO
	mapper.UpdatePlayerFromDTO(player, playerUpdate)

	// Persist changes
	if err := s.PlayerRepository.UpdatePlayer(ctx, playerID, player); err != nil {
		return nil, fmt.Errorf("failed to update player: %w", err)
	}

	// Update team associations if provided
	if len(playerUpdate.TeamIDs) > 0 {
		if err := s.updatePlayerTeams(ctx, player.ID, playerUpdate.TeamIDs); err != nil {
			return nil, err
		}
	}

	return player, nil
}

func (s *playerService) DeletePlayer(ctx context.Context, id uint64) error {
	return s.PlayerRepository.DeletePlayer(ctx, id)
}

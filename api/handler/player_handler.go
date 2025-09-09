package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/adapter/mapper"
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"github.com/EdwinRincon/browersfc-api/helper"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	PlayerDomainService *domainservice.PlayerDomainService
	PlayerMapper        *mapper.PlayerMapper
}

func NewPlayerHandler(playerDomainService *domainservice.PlayerDomainService) *PlayerHandler {
	return &PlayerHandler{
		PlayerDomainService: playerDomainService,
		PlayerMapper:        mapper.NewPlayerMapper(),
	}
}

// CreatePlayer godoc
// @Summary      Create a new player
// @Description  Creates a new player with the provided data
// @Tags         players
// @ID           createPlayer
// @Accept       json
// @Produce      json
// @Param        player  body      dto.CreatePlayerRequest  true  "Player data"
// @Success      201     {object}  dto.PlayerShort  "Player created successfully"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      409     {object}  helper.AppError "Conflict (e.g., nickname exists)"
// @Failure      500     {object}  helper.AppError "Internal server error"
// @Router       /admin/players [post]
// @Security     BearerAuth
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var createRequest dto.CreatePlayerRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid player data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Convert DTO to domain entity
	player := h.PlayerMapper.DTOToDomain(&createRequest)

	// Create player through domain service
	err := h.PlayerDomainService.CreatePlayer(ctx, player)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordAlreadyExists):
			helper.WriteErrorResponse(c, helper.NewConflictError("nick_name", "Nickname already exists"))
			return
		case errors.Is(err, constants.ErrInvalidData):
			helper.WriteErrorResponse(c, helper.NewBadRequestError("player", "Invalid player data"))
			return
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	// Return short player response
	playerShort := h.PlayerMapper.DomainToShortDTO(player)
	helper.WriteSuccessResponse(c, http.StatusCreated, playerShort, "Player created successfully")
}

// GetPlayerByID godoc
// @Summary      Get a player by ID
// @Description  Returns the details of a player by its ID
// @Tags         players
// @ID           getPlayerByID
// @Param        id   path      int  true  "Player ID"
// @Success      200  {object}  dto.PlayerResponse  "Player retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Player not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /players/{id} [get]
// @Security     BearerAuth
func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	playerID := c.Param("id")
	id, err := strconv.ParseUint(playerID, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid player ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	player, err := h.PlayerDomainService.GetPlayerByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("player"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	playerResponse := h.PlayerMapper.DomainToDTO(player)

	helper.WriteSuccessResponse(c, http.StatusOK, playerResponse, "Player retrieved successfully")
}

// GetPlayerByNickName godoc
// @Summary      Get a player by nickname
// @Description  Retrieves a player by their nickname
// @Tags         players
// @ID           getPlayerByNickName
// @Param        nickname  path      string  true  "Nickname"
// @Success      200       {object}  dto.PlayerResponse "Player retrieved successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      404       {object}  helper.AppError "Player not found"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /players/nickname/{nickname} [get]
// @Security     BearerAuth
func (h *PlayerHandler) GetPlayerByNickName(c *gin.Context) {
	nickname := c.Param("nickname")

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	player, err := h.PlayerDomainService.GetPlayerByNickName(ctx, nickname)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("player"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	playerResponse := h.PlayerMapper.DomainToDTO(player)

	helper.WriteSuccessResponse(c, http.StatusOK, playerResponse, "Player retrieved successfully")
}

// GetPaginatedPlayers godoc
// @Summary Get paginated players
// @Description Retrieve players with pagination, optional sorting, and ordering.
// @Tags players
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Param pageSize query int false "Page size" default(10)
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order" Enums(asc, desc) default(desc)
// @Success 200 {object} helper.AppSuccess{data=helper.PaginatedResponse{items=[]dto.PlayerResponse, totalCount=int}}
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /players [get]
// @Security BearerAuth
func (h *PlayerHandler) GetPaginatedPlayers(c *gin.Context) {
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 0 {
		page = 0
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 10
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Validate sort field
	if err := helper.ValidateSort(model.Player{}, sort); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	players, total, err := h.PlayerDomainService.GetPaginatedPlayers(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      h.PlayerMapper.DomainListToDTO(players),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Players retrieved successfully")
}

// UpdatePlayer godoc
// @Summary      Update an existing player
// @Description  Updates an existing player's information by ID
// @Tags         players
// @ID           updatePlayer
// @Accept       json
// @Produce      json
// @Param        id      path      int           true  "Player ID"
// @Param        player  body      dto.UpdatePlayerRequest true  "Updated player data"
// @Success      200     {object}  dto.PlayerShort  "Player updated successfully"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      404     {object}  helper.AppError "Player not found"
// @Failure      409     {object}  helper.AppError "Conflict (e.g., nickname exists)"
// @Failure      500     {object}  helper.AppError "Internal server error"
// @Router       /admin/players/{id} [put]
// @Security     BearerAuth
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid player ID"))
		return
	}

	var playerUpdateDTO dto.UpdatePlayerRequest
	if err = c.ShouldBindJSON(&playerUpdateDTO); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid player data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Convert DTO to domain entity
	playerUpdate := h.PlayerMapper.UpdateDTOToDomain(&playerUpdateDTO)

	updatedPlayer, err := h.PlayerDomainService.UpdatePlayer(ctx, id, playerUpdate)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("player"))
		} else if errors.Is(err, constants.ErrRecordAlreadyExists) {
			helper.WriteErrorResponse(c, helper.NewConflictError("nick_name", "Nickname already exists"))
		} else if errors.Is(err, constants.ErrInvalidData) {
			helper.WriteErrorResponse(c, helper.NewBadRequestError("player", "Invalid player data"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	response := h.PlayerMapper.DomainToShortDTO(updatedPlayer)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Player updated successfully")
}

// DeletePlayer godoc
// @Summary      Delete a player
// @Description  Deletes a player by its ID
// @Tags         players
// @ID           deletePlayer
// @Param        id   path      int  true  "Player ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Player not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /admin/players/{id} [delete]
// @Security     BearerAuth
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid player ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = h.PlayerDomainService.DeletePlayer(ctx, id)
	if err != nil && !errors.Is(err, constants.ErrRecordNotFound) {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

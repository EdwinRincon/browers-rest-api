package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	PlayerService service.PlayerService
}

func NewPlayerHandler(playerService service.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		PlayerService: playerService,
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
// @Security     ApiKeyAuth
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var createRequest dto.CreatePlayerRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.RespondWithError(c, helper.ProcessValidationError(err, "body", "Invalid player data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	createdPlayer, err := h.PlayerService.CreatePlayer(ctx, &createRequest)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordAlreadyExists):
			helper.RespondWithError(c, helper.Conflict("nick_name", "Nickname already exists"))
			return
		default:
			helper.RespondWithError(c, helper.InternalError(err))
			return
		}
	}

	helper.HandleSuccess(c, http.StatusCreated, createdPlayer, "Player created successfully")
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
// @Security     ApiKeyAuth
func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	playerID := c.Param("id")
	id, err := strconv.ParseUint(playerID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid player ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	player, err := h.PlayerService.GetPlayerByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("player"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	playerResponse := mapper.ToPlayerResponse(player)

	helper.HandleSuccess(c, http.StatusOK, playerResponse, "Player retrieved successfully")
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
// @Security     ApiKeyAuth
func (h *PlayerHandler) GetPlayerByNickName(c *gin.Context) {
	nickname := c.Param("nickname")

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	player, err := h.PlayerService.GetPlayerByNickName(ctx, nickname)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("player"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	playerResponse := mapper.ToPlayerResponse(player)

	helper.HandleSuccess(c, http.StatusOK, playerResponse, "Player retrieved successfully")
}

// GetPaginatedPlayers godoc
// @Summary      Get paginated players
// @Description  Retrieves a paginated list of players with sorting and ordering
// @Tags         players
// @ID           getPaginatedPlayers
// @Param        page      query     int     false  "Page number (0-based)"
// @Param        pageSize  query     int     false  "Number of items per page (default 10)"
// @Param        sort      query     string  false  "Sort field (e.g., nickname, rating)"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200       {object}  map[string]interface{} "Players retrieved successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /players [get]
// @Security     ApiKeyAuth
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
		helper.RespondWithError(c, helper.BadRequest("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	players, total, err := h.PlayerService.GetPaginatedPlayers(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToPlayerResponseList(players),
		TotalCount: total,
	}

	helper.HandleSuccess(c, http.StatusOK, response, "Players retrieved successfully")
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
// @Security     ApiKeyAuth
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid player ID"))
		return
	}

	var playerUpdateDTO dto.UpdatePlayerRequest
	if err = c.ShouldBindJSON(&playerUpdateDTO); err != nil {
		helper.RespondWithError(c, helper.ProcessValidationError(err, "body", "Invalid player data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedPlayer, err := h.PlayerService.UpdatePlayer(ctx, &playerUpdateDTO, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("player"))
		} else if errors.Is(err, constants.ErrRecordAlreadyExists) {
			helper.RespondWithError(c, helper.Conflict("nick_name", "Nickname already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	response := mapper.ToPlayerShort(updatedPlayer)
	helper.HandleSuccess(c, http.StatusOK, response, "Player updated successfully")
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
// @Security     ApiKeyAuth
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid player ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = h.PlayerService.DeletePlayer(ctx, id)
	if err != nil && !errors.Is(err, constants.ErrRecordNotFound) {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

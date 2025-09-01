package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
)

type PlayerTeamHandler struct {
	PlayerTeamService service.PlayerTeamService
}

func NewPlayerTeamHandler(playerTeamService service.PlayerTeamService) *PlayerTeamHandler {
	return &PlayerTeamHandler{
		PlayerTeamService: playerTeamService,
	}
}

// CreatePlayerTeam godoc
// @Summary      Create a new player-team relationship
// @Description  Creates a new player-team relationship with the provided data
// @Tags         playerTeams
// @ID           createPlayerTeam
// @Accept       json
// @Produce      json
// @Param        playerTeam  body      dto.CreatePlayerTeamRequest  true  "Player-Team relationship data"
// @Success      201   {object}  dto.PlayerTeamResponse  "Player-Team relationship created successfully"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      409   {object}  helper.AppError "Conflict (e.g., date overlap)"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /admin/player-teams [post]
// @Security     ApiKeyAuth
func (h *PlayerTeamHandler) CreatePlayerTeam(c *gin.Context) {
	var createRequest dto.CreatePlayerTeamRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.RespondWithError(c, helper.ProcessValidationError(err, "body", "Invalid player-team data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeamResponse, err := h.PlayerTeamService.CreatePlayerTeam(ctx, &createRequest)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrPlayerNotFound):
			helper.RespondWithError(c, helper.NotFound("player"))
			return
		case errors.Is(err, constants.ErrTeamNotFound):
			helper.RespondWithError(c, helper.NotFound("team"))
			return
		case errors.Is(err, constants.ErrSeasonNotFound):
			helper.RespondWithError(c, helper.NotFound("season"))
			return
		case errors.Is(err, constants.ErrOverlappingDates):
			helper.RespondWithError(c, helper.Conflict("date_range", "Date range overlaps with an existing record"))
			return
		default:
			helper.RespondWithError(c, helper.InternalError(err))
			return
		}
	}

	helper.HandleSuccess(c, http.StatusCreated, playerTeamResponse, "Player-team relationship created successfully")
}

// GetPlayerTeamByID godoc
// @Summary      Get player-team relationship by ID
// @Description  Retrieves a player-team relationship by its ID
// @Tags         playerTeams
// @ID           getPlayerTeamByID
// @Param        id  path      int  true  "PlayerTeam ID"
// @Success      200  {object}  dto.PlayerTeamResponse  "Player-team relationship retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /player-teams/{id} [get]
// @Security     ApiKeyAuth
func (h *PlayerTeamHandler) GetPlayerTeamByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeamResponse, err := h.PlayerTeamService.GetPlayerTeamByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("player-team relationship"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, playerTeamResponse, "Player-team relationship retrieved successfully")
}

// GetPlayerTeamsByPlayerID godoc
// @Summary      Get player-team relationships by player ID
// @Description  Retrieves all team relationships for a specific player
// @Tags         playerTeams
// @ID           get-player-teams-by-player-id
// @Param        id  path      int  true  "Player ID"
// @Produce      json
// @Success      200  {array}   dto.PlayerTeamResponse  "Player-team relationships retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Player not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /players/{id}/teams [get]
// @Security     ApiKeyAuth
func (h *PlayerTeamHandler) GetPlayerTeamsByPlayerID(c *gin.Context) {
	playerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid player ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeamResponses, err := h.PlayerTeamService.GetPlayerTeamsByPlayerID(ctx, playerID)
	if err != nil {
		if errors.Is(err, constants.ErrPlayerNotFound) {
			helper.RespondWithError(c, helper.NotFound("player"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, playerTeamResponses, "Player-team relationships retrieved successfully")
}

// GetPlayerTeamsByTeamID godoc
// @Summary      Get player-team relationships by team ID
// @Description  Retrieves all player relationships for a specific team
// @Tags         playerTeams
// @ID           getPlayerTeamsByTeamID
// @Param        id  path      int  true  "Team ID"
// @Produce      json
// @Success      200  {array}   dto.PlayerTeamResponse  "Player-team relationships retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Team not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /teams/{id}/players [get]
// @Security     ApiKeyAuth
func (h *PlayerTeamHandler) GetPlayerTeamsByTeamID(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid team ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeamResponses, err := h.PlayerTeamService.GetPlayerTeamsByTeamID(ctx, teamID)
	if err != nil {
		if errors.Is(err, constants.ErrTeamNotFound) {
			helper.RespondWithError(c, helper.NotFound("team"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, playerTeamResponses, "Player-team relationships retrieved successfully")
}

// GetPlayerTeamsBySeasonID godoc
// @Summary      Get player-team relationships by season ID
// @Description  Retrieves all player-team relationships for a specific season
// @Tags         playerTeams
// @ID           getPlayerTeamsBySeasonID
// @Param        id  path      int  true  "Season ID"
// @Produce      json
// @Success      200  {array}   dto.PlayerTeamResponse  "Player-team relationships retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Season not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /seasons/{id}/player-teams [get]
// @Security     ApiKeyAuth
func (h *PlayerTeamHandler) GetPlayerTeamsBySeasonID(c *gin.Context) {
	seasonID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid season ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeamResponses, err := h.PlayerTeamService.GetPlayerTeamsBySeasonID(ctx, seasonID)
	if err != nil {
		if errors.Is(err, constants.ErrSeasonNotFound) {
			helper.RespondWithError(c, helper.NotFound("season"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, playerTeamResponses, "Player-team relationships retrieved successfully")
}

// GetPaginatedPlayerTeams godoc
// @Summary      Get paginated player-team relationships
// @Description  Retrieves a paginated list of player-team relationships with sorting and ordering
// @Tags         playerTeams
// @ID           getPaginatedPlayerTeams
// @Param        page      query     int     false  "Page number (0-based)"
// @Param        pageSize  query     int     false  "Number of items per page (default 10)"
// @Param        sort      query     string  false  "Sort field (e.g., created_at)"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200       {object}  map[string]interface{} "Player-team relationships retrieved successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /player-teams [get]
// @Security     ApiKeyAuth
func (h *PlayerTeamHandler) GetPaginatedPlayerTeams(c *gin.Context) {
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
	if err := helper.ValidateSort(model.PlayerTeam{}, sort); err != nil {
		helper.RespondWithError(c, helper.BadRequest("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeams, total, err := h.PlayerTeamService.GetPaginatedPlayerTeams(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      playerTeams,
		TotalCount: total,
	}

	helper.HandleSuccess(c, http.StatusOK, response, "Player-team relationships retrieved successfully")
}

// UpdatePlayerTeam godoc
// @Summary      Update a player-team relationship
// @Description  Updates an existing player-team relationship by ID
// @Tags         playerTeams
// @ID           updatePlayerTeam
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "PlayerTeam ID"
// @Param        playerTeam  body      dto.UpdatePlayerTeamRequest  true  "Updated player-team data"
// @Success      200  {object}  dto.PlayerTeamResponse  "Player-team relationship updated successfully"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Not found"
// @Failure      409  {object}  helper.AppError "Conflict (e.g., date overlap)"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /admin/player-teams/{id} [put]
// @Security     ApiKeyAuth
func (h *PlayerTeamHandler) UpdatePlayerTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid ID"))
		return
	}

	var updateRequest dto.UpdatePlayerTeamRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		helper.RespondWithError(c, helper.ProcessValidationError(err, "body", "Invalid player-team data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeamResponse, err := h.PlayerTeamService.UpdatePlayerTeam(ctx, id, &updateRequest)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordNotFound):
			helper.RespondWithError(c, helper.NotFound("player-team relationship"))
			return
		case errors.Is(err, constants.ErrOverlappingDates):
			helper.RespondWithError(c, helper.Conflict("date_range", "Date range overlaps with an existing record"))
			return
		default:
			helper.RespondWithError(c, helper.InternalError(err))
			return
		}
	}

	helper.HandleSuccess(c, http.StatusOK, playerTeamResponse, "Player-team relationship updated successfully")
}

// DeletePlayerTeam godoc
// @Summary      Delete a player-team relationship
// @Description  Deletes a player-team relationship by ID
// @Tags         playerTeams
// @ID           deletePlayerTeam
// @Param        id  path      int  true  "PlayerTeam ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /admin/player-teams/{id} [delete]
// @Security     ApiKeyAuth
func (h *PlayerTeamHandler) DeletePlayerTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.PlayerTeamService.DeletePlayerTeam(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("player-team relationship"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

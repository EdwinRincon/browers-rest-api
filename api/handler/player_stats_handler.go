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

type PlayerStatsHandler struct {
	PlayerStatsService service.PlayerStatsService
}

func NewPlayerStatsHandler(playerStatsService service.PlayerStatsService) *PlayerStatsHandler {
	return &PlayerStatsHandler{
		PlayerStatsService: playerStatsService,
	}
}

// CreatePlayerStat godoc
// @Summary      Create a new player statistic
// @Description  Creates a new player statistic with the provided data
// @Tags         player-stats
// @ID           createPlayerStat
// @Accept       json
// @Produce      json
// @Param        playerStat  body      dto.CreatePlayerStatRequest  true  "Player Statistic data"
// @Success      201   {object}  dto.PlayerStatResponse  "Player Statistic created successfully"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      404   {object}  helper.AppError "Related entity not found"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /admin/player-stats [post]
// @Security     ApiKeyAuth
func (h *PlayerStatsHandler) CreatePlayerStat(c *gin.Context) {
	var createRequest dto.CreatePlayerStatRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", constants.MsgInvalidData))
		return
	}

	playerStat, err := h.PlayerStatsService.CreatePlayerStat(c.Request.Context(), &createRequest)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrPlayerNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgInvalidPlayerID))
			return
		case errors.Is(err, constants.ErrTeamNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgInvalidTeamID))
			return
		case errors.Is(err, constants.ErrSeasonNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgInvalidSeasonData))
			return
		case errors.Is(err, constants.ErrRecordNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgNotFound))
			return
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	response := mapper.ToPlayerStatResponse(playerStat)
	helper.WriteSuccessResponse(c, http.StatusCreated, response, "Player statistic created successfully")
}

// GetPlayerStatByID godoc
// @Summary      Get a player statistic by ID
// @Description  Returns a player statistic by its ID
// @Tags         player-stats
// @ID           getPlayerStatByID
// @Produce      json
// @Param        id  path      string  true  "Player Statistic ID"
// @Success      200  {object}  dto.PlayerStatResponse  "Player statistic retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid ID format"
// @Failure      404  {object}  helper.AppError "Player statistic not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /player-stats/{id} [get]
func (h *PlayerStatsHandler) GetPlayerStatByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidID))
		return
	}

	playerStat, err := h.PlayerStatsService.GetPlayerStatByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgNotFound))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	response := mapper.ToPlayerStatResponse(playerStat)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Player statistic retrieved successfully")
}

// GetPlayerStatsByPlayerID godoc
// @Summary      Get player statistics for a player
// @Description  Returns all player statistics for a specific player
// @Tags         player-stats
// @ID           getPlayerStatsByPlayerID
// @Produce      json
// @Param        id  path      string  true  "Player ID"
// @Success      200  {array}   dto.PlayerStatResponse  "Player statistics retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid player ID format"
// @Failure      404  {object}  helper.AppError "Player not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /players/{id}/stats [get]
func (h *PlayerStatsHandler) GetPlayerStatsByPlayerID(c *gin.Context) {
	playerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidPlayerID))
		return
	}

	playerStats, err := h.PlayerStatsService.GetPlayerStatsByPlayerID(c.Request.Context(), playerID)
	if err != nil {
		if errors.Is(err, constants.ErrPlayerNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgInvalidPlayerID))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	response := mapper.ToPlayerStatResponseList(playerStats)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Player statistics retrieved successfully")
}

// GetPlayerStatsByMatchID godoc
// @Summary      Get player statistics for a match
// @Description  Returns all player statistics for a specific match
// @Tags         player-stats
// @ID           getPlayerStatsByMatchID
// @Produce      json
// @Param        id  path      string  true  "Match ID"
// @Success      200  {array}   dto.PlayerStatResponse  "Player statistics retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid match ID format"
// @Failure      404  {object}  helper.AppError "Match not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /matches/{id}/stats [get]
func (h *PlayerStatsHandler) GetPlayerStatsByMatchID(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidMatchID))
		return
	}

	playerStats, err := h.PlayerStatsService.GetPlayerStatsByMatchID(c.Request.Context(), matchID)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgInvalidMatchID))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	response := mapper.ToPlayerStatResponseList(playerStats)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Player statistics retrieved successfully")
}

// GetPlayerStatsBySeasonID godoc
// @Summary      Get player statistics for a season
// @Description  Returns all player statistics for a specific season
// @Tags         player-stats
// @ID           getPlayerStatsBySeasonID
// @Produce      json
// @Param        id  path      string  true  "Season ID"
// @Success      200  {array}   dto.PlayerStatResponse  "Player statistics retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid season ID format"
// @Failure      404  {object}  helper.AppError "Season not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /seasons/{id}/stats [get]
func (h *PlayerStatsHandler) GetPlayerStatsBySeasonID(c *gin.Context) {
	seasonID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidSeasonData))
		return
	}

	playerStats, err := h.PlayerStatsService.GetPlayerStatsBySeasonID(c.Request.Context(), seasonID)
	if err != nil {
		if errors.Is(err, constants.ErrSeasonNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgInvalidSeasonData))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	response := mapper.ToPlayerStatResponseList(playerStats)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Player statistics retrieved successfully")
}

// GetPaginatedPlayerStats godoc
// @Summary      Get paginated player statistics
// @Description  Returns a paginated list of player statistics
// @Tags         player-stats
// @ID           getPaginatedPlayerStats
// @Produce      json
// @Param        page      query     int     false  "Page number (0-indexed)"
// @Param        pageSize  query     int     false  "Page size"
// @Param        sort      query     string  false  "Sort field"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200  {object}  helper.PaginatedResponse  "Player statistics retrieved successfully"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /player-stats [get]
func (h *PlayerStatsHandler) GetPaginatedPlayerStats(c *gin.Context) {
	sort := c.DefaultQuery("sort", "id")
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
	if err := helper.ValidateSort(model.PlayerStat{}, sort); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerStats, total, err := h.PlayerStatsService.GetPaginatedPlayerStats(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToPlayerStatResponseList(playerStats),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Player statistics retrieved successfully")
}

// UpdatePlayerStat godoc
// @Summary      Update a player statistic
// @Description  Updates a player statistic with the provided data
// @Tags         player-stats
// @ID           updatePlayerStat
// @Accept       json
// @Produce      json
// @Param        id          path      string                     true  "Player Statistic ID"
// @Param        playerStat  body      dto.UpdatePlayerStatRequest  true  "Player Statistic update data"
// @Success      200  {object}  dto.PlayerStatResponse  "Player statistic updated successfully"
// @Failure      400  {object}  helper.AppError "Invalid input or ID format"
// @Failure      404  {object}  helper.AppError "Player statistic not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /admin/player-stats/{id} [put]
// @Security     ApiKeyAuth
func (h *PlayerStatsHandler) UpdatePlayerStat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidID))
		return
	}

	var updateRequest dto.UpdatePlayerStatRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", constants.MsgInvalidData))
		return
	}

	playerStat, err := h.PlayerStatsService.UpdatePlayerStat(c.Request.Context(), id, &updateRequest)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgNotFound))
		case errors.Is(err, constants.ErrTeamNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgInvalidTeamID))
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	response := mapper.ToPlayerStatResponse(playerStat)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Player statistic updated successfully")
}

// DeletePlayerStat godoc
// @Summary      Delete a player statistic
// @Description  Deletes a player statistic by ID
// @Tags         player-stats
// @ID           deletePlayerStat
// @Produce      json
// @Param        id  path      string  true  "Player Statistic ID"
// @Success      200  {object}  helper.AppSuccess  "Player statistic deleted successfully"
// @Failure      400  {object}  helper.AppError "Invalid ID format"
// @Failure      404  {object}  helper.AppError "Player statistic not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /admin/player-stats/{id} [delete]
// @Security     ApiKeyAuth
func (h *PlayerStatsHandler) DeletePlayerStat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidID))
		return
	}

	err = h.PlayerStatsService.DeletePlayerStat(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError(constants.MsgNotFound))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	helper.WriteSuccessResponse(c, http.StatusOK, nil, "Player statistic deleted successfully")
}

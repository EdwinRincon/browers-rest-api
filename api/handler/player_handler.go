package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
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

func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var player model.Players
	if err := c.ShouldBindJSON(&player); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err := h.PlayerService.CreatePlayer(ctx, &player)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create player", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, player, "Player created successfully")
}

func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidPlayerID, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	player, err := h.PlayerService.GetPlayerByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPlayerNotFound) {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrPlayerNotFound.Error(), ""), false)
			return
		}
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to retrieve player", err.Error()), false)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, player, "Player retrieved successfully")
}

func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	players, err := h.PlayerService.ListPlayers(ctx, page)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list players", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, players, "Players listed successfully")
}

func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidPlayerID, err.Error()), true)
		return
	}

	var player model.Players
	if err := c.ShouldBindJSON(&player); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}
	player.ID = id

	ctx := c.Request.Context()
	err = h.PlayerService.UpdatePlayer(ctx, &player)
	if err != nil {
		if errors.Is(err, repository.ErrPlayerNotFound) {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrPlayerNotFound.Error(), ""), true)
			return
		}
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update player", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, player, "Player updated successfully")
}

func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidPlayerID, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.PlayerService.DeletePlayer(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete player", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Player deleted successfully")
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/constants"
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

func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var player model.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.PlayerService.CreatePlayer(ctx, &player)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, player, "Player created successfully")
}

func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	playerID := c.Param("id")
	id, err := strconv.ParseUint(playerID, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	player, err := h.PlayerService.GetPlayerByID(ctx, id)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, player, "Player found successfully")
}

func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	playerID := c.Param("id")
	id, err := strconv.ParseUint(playerID, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	var player model.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	player.ID = id
	err = h.PlayerService.UpdatePlayer(ctx, &player)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, player, "Player updated successfully")
}

func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	playerID := c.Param("id")
	id, err := strconv.ParseUint(playerID, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.PlayerService.DeletePlayer(ctx, id)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Player deleted successfully")
}

func (h *PlayerHandler) GetAllPlayers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}
	ctx := c.Request.Context()
	players, err := h.PlayerService.GetAllPlayers(ctx, page)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, players, "Player retrieved successfully")
}
